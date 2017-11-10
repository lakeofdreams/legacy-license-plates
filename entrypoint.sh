#!/bin/bash
########### SIGTERM handler ############
function _term() {
   echo "Stopping container."
   echo "SIGTERM received, shutting down the server!"
   ${DOMAIN_HOME}/bin/stopWebLogic.sh
}
########### SIGKILL handler ############
function _kill() {
   echo "SIGKILL received, shutting down the server!"
   kill -9 $childPID
}
# Set SIGTERM handler
trap _term SIGTERM
# Set SIGKILL handler
trap _kill SIGKILL
# wait for DB to come up
echo "Waiting for database to be available..."
while ! nc database 1521 < /dev/null; do sleep 9; done
# sleep for some extra time for db setup
sleep 80
# if there are any ddl files, run them
if [ -e /u01/oracle/create-users.sql ];
then
   echo "Creating users..."
   java -classpath /u01/oracle/wlserver/server/lib/weblogic.jar \
   utils.Schema \
   jdbc:oracle:thin:@//database:1521/OraDoc.my.domain.com \
   oracle.jdbc.OracleDriver \
   -u "sys as sysdba" \
   -p "Welcome1" \
   -verbose \  
   /u01/oracle/create-users.sql
fi
# userid hard coded for now - fix this later
if [ -e /u01/oracle/create-schema.sql ];
then
   echo "Creating schema..."
   java -classpath /u01/oracle/wlserver/server/lib/weblogic.jar \
   utils.Schema \
   jdbc:oracle:thin:@//database:1521/OraDoc.my.domain.com \
   oracle.jdbc.OracleDriver \
   -u "c##mark" \
   -p "Welcome1" \
   -verbose \
   /u01/oracle/create-schema.sql
fi


DOMAIN_HOME=/u01/oracle/domains/base_domain

# Start Admin Server and tail the logs
echo "Starting WebLogic..."
${DOMAIN_HOME}/startWebLogic.sh &
childPID=$!
wait $childPID
