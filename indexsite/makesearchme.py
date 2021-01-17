#!/usr/bin/python

import sys
import os
import MySQLdb

if len(sys.argv) > 1:
   domain = sys.argv[1]
else:
   print "Enter the domain:",
   domain = raw_input()

base_dir = "/var/www/"

searchme_loc = base_dir + domain + "/searchme.html"

print(os.environ)

user = os.environ['DBUSER']
passwd = os.environ['DBPASS']
db = os.environ['DBNAME']

conn = MySQLdb.connect (host = "localhost",
                           user = user,
                           passwd = passwd,
                           db = db)
cursor = conn.cursor()

cursor.execute("SELECT artist_id FROM cdartist")


searchmefile = open(searchme_loc, 'w')
searchmefile.write("<html><head></head><body>\n")


row = cursor.fetchone()
while row is not None:
   myartist_id = row[0]
   link_line = "<a href=\"/cdcoll.php?artist_id=" + str(myartist_id) + "\"></a>\n"
   searchmefile.write(link_line)
   row = cursor.fetchone()

link_line = "<a href=\"/\"></a>\n"
searchmefile.write(link_line)
link_line = "<a href=\"/searchsite.php\"></a>\n"
searchmefile.write(link_line)
searchmefile.write("</body></html>\n")
cursor.close()
conn.close()
searchmefile.close()

