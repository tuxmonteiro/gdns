#!/bin/bash

DOMAIN=$1
TOKEN=$GDNS_TOKEN
CONTACT="root.${DOMAIN}."
GDNS=$GDNS_HOST
CHALLENGE="$(dd if=/dev/urandom bs=1 count=200 2> /dev/null | base64 -b300 | sed -e 's/\+/_/g' -e 's/[= ]//g' -e 's/\//_/g' | cut -c1-200)"
CURL="curl -v -HContent-type:application/json -HX-Auth-Token:${TOKEN}"

echo ${CURL} -X POST -d '{"domain": {"name": "'${DOMAIN}'", "type": "MASTER", "ttl": "86400", "notes": "A domain", "primary_ns": "ns1.'${DOMAIN}'.", "contact": "'${CONTACT}'", "refresh": 10800, "retry": 3600, "expire": 604800, "minimum": 10800,"authority_type": "M"}}' ${GDNS}/domains.json
DATA="{\"domain\": {\"name\": \"${DOMAIN}\", \"type\": \"MASTER\", \"ttl\": \"86400\", \"notes\": \"A domain\", \"primary_ns\": \"ns1.${DOMAIN}.\", \"contact\": \"${CONTACT}\", \"refresh\": 10800, \"retry\": 3600, \"expire\": 604800, \"minimum\": 10800,\"authority_type\": \"M\"}}"
${CURL} -X POST -d "$DATA" ${GDNS}/domains.json | jq -r .domain.id | tee domain.id | sed 's/^/Domain ID: /'
DID=$(cat domain.id)

echo ${CURL} -X POST -d '{"record": {"name": "@", "type": "NS", "content":"ns1.'${DOMAIN}'."}}' ${GDNS}/domains/${DID}/records.json
${CURL} -X POST -d '{"record": {"name": "@", "type": "NS", "content":"ns1.'${DOMAIN}'."}}' ${GDNS}/domains/${DID}/records.json | jq .
echo ${CURL} -X POST -d '{"record": {"name": "@", "type": "NS", "content":"ns2.'${DOMAIN}'."}}' ${GDNS}/domains/${DID}/records.json
${CURL} -X POST -d '{"record": {"name": "@", "type": "NS", "content":"ns2.'${DOMAIN}'."}}' ${GDNS}/domains/${DID}/records.json | jq .
echo ${CURL} -X POST -d '{"record": {"name": "@", "type": "NS", "content":"ns3.'${DOMAIN}'."}}' ${GDNS}/domains/${DID}/records.json
${CURL} -X POST -d '{"record": {"name": "@", "type": "NS", "content":"ns3.'${DOMAIN}'."}}' ${GDNS}/domains/${DID}/records.json | jq .
echo ${CURL} -X POST -d '{"record": {"name": "ns1.'${DOMAIN}'.", "type": "A", "content":"10.236.30.95"}}' ${GDNS}/domains/${DID}/records.json
${CURL} -X POST -d '{"record": {"name": "ns1.'${DOMAIN}'.", "type": "A", "content":"10.236.30.95"}}' ${GDNS}/domains/${DID}/records.json | jq .
echo ${CURL} -X POST -d '{"record": {"name": "ns2.'${DOMAIN}'.", "type": "A", "content":"10.236.30.96"}}' ${GDNS}/domains/${DID}/records.json
${CURL} -X POST -d '{"record": {"name": "ns2.'${DOMAIN}'.", "type": "A", "content":"10.236.30.96"}}' ${GDNS}/domains/${DID}/records.json | jq .
echo ${CURL} -X POST -d '{"record": {"name": "ns3.'${DOMAIN}'.", "type": "A", "content":"10.236.30.103"}}' ${GDNS}/domains/${DID}/records.json
${CURL} -X POST -d '{"record": {"name": "ns3.'${DOMAIN}'.", "type": "A", "content":"10.236.30.103"}}' ${GDNS}/domains/${DID}/records.json | jq .

DATA="{\"record\": {\"name\": \"_acme-challenge.${DOMAIN}.\", \"type\": \"TXT\", \"content\": \"\\\"${CHALLENGE}\\\"\"}}"
echo "${CURL} -X POST -d \"$DATA\" ${GDNS}/domains/${DID}/records.json"
${CURL} -X POST -d "${DATA}" ${GDNS}/domains/${DID}/records.json | jq .

echo ${CURL} -X POST ${GDNS}/bind9/export.json
${CURL} -X POST ${GDNS}/bind9/export.json | jq .
echo ${CURL} -X POST ${GDNS}/bind9/schedule_export.json
${CURL} -X POST ${GDNS}/bind9/schedule_export.json | jq .

# EOF
