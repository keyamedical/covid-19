http = require('https');

const host = 'covid-detc.us.keyayun.net';
const token = 'the bearer token';

newWorkItem = (affectedSOPInstanceUID, studyInstanceUID) => new Promise((resolve, reject) => {
  options = {
    method: 'POST',
    hostname: host,
    path: `/workitems?${affectedSOPInstanceUID}`,
    headers: {
      Authorization: `Bearer ${token}`
    }
  };
  req = http.request(options, res => {
    chunks = [];
    res.on('data',chunk => {
      chunks.push(chunk);
    });
    res.on('end',() => {
      body = Buffer.concat(chunks);
      if (res.statusCode !== 201) {
        reject(new Error(body.toString()));
      }
      resolve();
    });
    res.on('error', (err) => {
      reject(err);
    })
  });
  req.write(`{\n    \"00741204\": \"COVID-19\",\n    \"00404021\": {\n        \"0020000D\": \"${studyInstanceUID}\"\n    }\n}`);
  req.end();
});

getWorkItem = affectedSOPInstanceUID => new Promise((resolve, reject) => {
  options = {
    method: 'GET',
    hostname: host,
    path: `/workitems/${affectedSOPInstanceUID}`,
    headers: {
      Authorization: `Bearer ${token}`
    }
  };
  req = http.request(options, res => {
    chunks = [];
    res.on('data',chunk => {
      chunks.push(chunk);
    });
    res.on('end',() => {
      body = Buffer.concat(chunks);
      if (res.statusCode !== 200) {
        reject(new Error(body.toString()));
      }
      resolve(body.toString());
    });
    res.on('error', (err) => {
      reject(err);
    })
  });
  req.end();
});

testSOPUID = '1.3.45.214';
testStudyUID = '2.16.840.1.113662.2.1.99999.5175439602988854';
newWorkItem(testSOPUID, testStudyUID)
  .then(() => {
    return getWorkItem(testSOPUID);
  })
  .then((res) => {
    console.log(res);
  })
  .catch(err => {
    console.error(err.stack)
  });

