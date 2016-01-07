import MySQLdb

from backends.settings import DATABASES as db

mysqldb = db['default']

dbname = mysqldb['NAME']
dbuser = mysqldb['USER']
dbpasswd = mysqldb['PASSWORD']
dbhost = mysqldb['HOST']
dbport = mysqldb['PORT']
