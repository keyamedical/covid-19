#!/bin/bash
set -e
url="https://covid-detc.us.keyayun.net"
token="the bearer token"

function newWorkItem() {
  affectedSOPInstanceUID=$1
  studyInstanceUID=$2

  status_code=`curl -w %{http_code} --output /dev/null -X POST \
    ''${url}'/workitems?'${affectedSOPInstanceUID}'' \
    -H 'Authorization: Bearer '${token}'' \
    -d '{
      "00741204": "COVID-19",
      "00404021": {
          "0020000D": "'${studyInstanceUID}'"
      }
  }'`
  if [[ "$status_code" -ne 201 ]] ; then
    echo "status code: ${status_code}"
    exit 1
  fi
}

function getWorkItem() {
  affectedSOPInstanceUID=$1
  curl -X GET \
    ''${url}'/workitems/'${affectedSOPInstanceUID}'' \
    -H 'Authorization: Bearer '${token}''
}

testSOPUID="1.3.45.214"
testStudyUID="2.16.840.1.113662.2.1.99999.5175439602988854"
newWorkItem $testSOPUID $testStudyUID
getWorkItem $testSOPUID
