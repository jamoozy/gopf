# Endpoints

GOPF has the following endpoints:

**Tagging a File**
```
GET /tag?t={tagname}&f={filename}
```
Tag a file with a given name.  Both `filename` and `tagname` must be the full names of the tag and file as defined in the database.  Returns either 

**Listing Tags**
```
GET /tags
```
Gets a list of tags
