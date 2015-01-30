#!/usr/bin/ruby -w
#
# Encapsulates the functions of the DB.

require 'sqlite3'
require 'digest'

# Class handling the DB interface.
class DBI
  attr_accessor :db_name

  def initialize(db_name='gopf.db')
    @db_name = db_name
  end

  # Create the DB with tables for file references, tags, and the relationship
  # between the two.
  def with_db
    SQLite3::Database.new(@db_name) do |db|
      # Create, if D.N.E. (previously).
      db.execute("create table 'files' (
                    id integer not null primary key autoincrement,
                    hash text not null,
                    path text not null
                  )") unless `sqlite3 #@db_name '.schema files'`.size > 0
      db.execute("create table 'tags' (
                    id integer not null primary key autoincrement,
                    name text not null unique
                  )") unless `sqlite3 #@db_name '.schema tags'`.size > 0
      db.execute("create table 'file_tags' (
                    file_id integer not null,
                    tag_id integer not null,
                    foreign key (file_id) references files(id),
                    foreign key (tag_id) references tags(id)
                  )") unless `sqlite3 #@db_name '.schema file_tags'`.size > 0

      # Perform whatever.
      yield db if block_given?
    end
  end

  # Adds a file to the DB.
  def add_file(path)
    with_db do |db|
      db.query('insert into files (path) values (?)', path).close
    end
  end

  # Ensure the path
  def ensure_file_in_db(path)
    with_db do |db|
      puts "Adding file \"#{path}\""
      hash = Digest::MD5.file(path).hexdigest
      db.query("select * from files where hash=? and path=?", [hash, path]) do |res|
        db.execute("insert into files (hash,path) values (?,?)", [hash, path]){|r|} unless res.any?
      end
    end
  end

  # Ensure the path
  def ensure_tag_in_db(name)
    with_db do |db|
      db.query("select * from tags where name=?", [name]) do |res|
        if res.any? then
          puts "Already have tag #{name}"
        else
          db.query("insert into tags (name) values (?)", [name]){|r|}
        end
      end
    end
  end

  # Moves a file from where it is to dir/
  def mv_file(path, dir)
    puts "Moving #{path} to #{dir}"
    with_db do |db|
      db.query('update files set path=? where path=?',
                 ["#{dir}/#{File.basename(path)}", path]){|r|}
      `mv #{path} #{dir}`
    end
  end

  # Tags a file.  The parameters can either be the DB IDs of the file/tag or
  # the path/name of the file/tag.
  def tag(file_noi, tag_noi)
    with_db do |db|
      puts "Tagging #{file_noi} as #{tag_noi}"

      begin
        db.query("select id from files where path=?", [file_noi]) {|res|
          file_noi=res.first[0]} unless file_noi == Fixnum
      rescue
        puts "File probably doesn't exist!"
        return
      end

      db.query("select id from tags where name=?", [tag_noi]) {|res|
        tag_noi=res.first[0]} unless tag_noi == Fixnum

      # Do the addition only if such a row doesn't already exist.
      db.query('select * from file_tags where file_id=? and tag_id=?', [file_noi, tag_noi]) do |res|
        if res.any? then
          puts "Already tagged."
        else
          db.execute('insert into file_tags (file_id,tag_id) values (?,?)', [file_noi, tag_noi]) do |r|
            puts "Added file_tag(#{file_noi},#{tag}) for \"#{file_noi}\" & \"#{tag}\""
          end
        end
      end
    end
  end
end


if __FILE__ == $0 then
  # Init DB interface.
  $dbi = DBI.new('gopf.db')

  # Make DB entry for each file
  Dir['*.mp4'].each do |f|
    $dbi.ensure_file_in_db(f)
  end if ARGV.include? '--files'

  # Convert each directory name to a tag, and convert each tag to a DB entry.
  #
  # This stanza is, effectively, deprecated as it was written for a very
  # specific purpose.  I'm keeping it here for posterity.
#  Dir['*'].reject{|f|!File.directory?(f)}.sort.each do |dir|
#    $dbi.ensure_tag_in_db dir
#
#    Dir["#{dir}/*.mp4"].reject{|f|!File.symlink?(f)}.each do |f|
#      target = File.basename(`readlink '#{f}'`).strip
#      $dbi.tag(target, dir)
#    end
#  end if ARGV.include? '--tags'

  # Move files to a data/ dir.
  Dir['*.mp4'].each do |f|
    $dbi.mv_file(f, 'data/')
  end if ARGV.include? '--data'

  puts "running from \"#{`pwd`.strip}\"" if ARGV.empty?
end
