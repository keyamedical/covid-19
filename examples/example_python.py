import requests

url = "https://covid-detc.us.keyayun.net"
token = "the bearer token"

def newWorkItem(affectedSOPInstanceUID, studyInstanceUID):
  apiUrl = url +  "/workitems?" + affectedSOPInstanceUID
  payload = "{\n    \"00741204\": \"COVID-19\",\n    \"00404021\": {\n        \"0020000D\": \""+ studyInstanceUID +"\"\n    }\n}"

  headers = {'Authorization': token}

  response = requests.request("POST", apiUrl, data=payload, headers=headers)

  if response.status_code != 201:
    raise Exception(response.text)

def getWorkItem(affectedSOPInstanceUID):
  apiUrl = url + "/workitems/" + affectedSOPInstanceUID

  headers = {'Authorization': token}

  response = requests.request("GET", apiUrl, headers=headers)
  if response.status_code != 200:
      raise Exception(response.text)
  print(response.json())

if __name__ == '__main__':
  testSOPUID = "1.3.45.214.7"
  testStudyUID = "2.16.840.1.113662.2.1.99999.5175439602988854"
  newWorkItem(testSOPUID, testStudyUID)
  getWorkItem(testSOPUID)
