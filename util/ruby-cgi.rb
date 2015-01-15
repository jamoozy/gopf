#!/usr/bin/env ruby

#  Copyright (c) 2011 Henrik Hodne
#
# Permission is hereby granted, free of charge, to any person
# obtaining a copy of this software and associated documentation
# files (the "Software"), to deal in the Software without
# restriction, including without limitation the rights to use,
# copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the
# Software is furnished to do so, subject to the following
# conditions:
#
# The above copyright notice and this permission notice shall be
# included in all copies or substantial portions of the Software.
#
# Except as contained in this notice, the name(s) of the above
# copyright holders shall not be used in advertising or otherwise
# to promote the sale, use or other dealings in this Software
# without prior written authorization.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
# EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
# OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
# NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
# HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
# WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
# FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
# OTHER DEALINGS IN THE SOFTWARE.

require 'cgi'
require 'erb'

$cgi = CGI.new
include ERB::Util

# Use this to output a header. This will output the header immediately,
# and ensures that it's printed before the other content
# This accepts either string(s) or a hash
# To force html content type, in your erb you must specify:
#   <% header "Content-Type" => "text/html" %>
#   or
#   name the filewith extension .html.erb
def header(*args)
  if args.length == 1
    args.first.each do |key, value|
      puts "#{key}: #{value}\r\n"
    end
  else
    args.each {|s| puts s+"\r\n" }
  end
  ($HEADERS ||= []).push(*args)
end

bind = binding # To prevent the erb script to have access to "f"

puts "\r\n" + File.open($cgi.path_translated, "r") {|f| ERB.new(f.read).result(bind) }
