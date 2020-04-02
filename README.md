
# COVID-19 Detection Service

<!-- TOC -->
- [1. Introduction](#1-introduction)
  - [1.1 Overview](#11-overview)
  - [1.2 Workflow](#12-workflow)
  - [1.3 CT Image Requirements](#13-ct-image-requirements)
- [2. RESTful API](#2-restful-api)
  - [2.1 Bearer Token](#21-bearer-token)
  - [2.2 HTTP Request/Response Payload Syntax](#22-http-requestresponse-payload-syntax)
  - [2.3 Request Analysis](#23-request-analysis)
    - [POST /workitems?{AffectedSOPInstanceUID}](#post-workitemsaffectedsopinstanceuid)
      - [AffectedSOPInstanceUID requirements](#affectedsopinstanceuid-requirements)
    - [Status codes](#status-codes)
    - [Required Keys](#required-keys)
    - [Request](#request)
    - [Response](#response)
  - [2.4 Query Status](#24-query-status)
    - [GET /workitems/{AffectedSOPInstanceUID}](#get-workitemsaffectedsopinstanceuid)
    - [Status codes](#status-codes-1)
    - [Request](#request-1)
    - [Response](#response-1)
    - [Cancellation](#cancellation)
    - [Completed](#completed)
- [3. DICOM Transfer](#3-dicom-transfer)
  - [Requirements](#requirments)
- [4. Free Trial Limitations](#4-free-trial-limitations)
<!-- /TOC -->

<div style="page-break-after: always;"></div>

## 1. Introduction

### 1.1 Overview

**COVID-19 Detection Service** is a SaaS application developed by KEYA MEDICAL LTD, aiming to detect COVID-19 disease based on chest CT images. The service can be integrated with workflow management system or worklist.

### 1.2 Workflow

1. Send chest CT Images to COVID-19 Server
2. Trigger analysis request to COVID-19 Server
3. *Wait for result...(~1 minute)*
4. Query status from COVID-19 Server

![workflow diagram](./docs/diagram-workflow.svg)

### 1.3 CT Image Requirements

1. Slice thickness **MUST** be less than or equal to 3mm
2. Number of Slices **MUST** be between 50 and 500

## 2. RESTful API

### 2.1 Bearer Token

All APIs **MUST** be called with **bearer token** included in HTTP request headers, for example

```http
POST /workitems?1.3.45.214 HTTP/1.1
Authorization: Bearer {token}
```

> The definition of bearer token **CAN** be found in [RFC 6750](https://tools.ietf.org/html/rfc6750#section-1.2)

### 2.2 HTTP Request/Response Payload Syntax

| Key | Description | Value
| :---: | --- | ---
| 0020000D | Study Instance UID
| **00400280** | **COVID-19 Prediction** | **Positive<br/>Negative**
| 00404005 | Scheduled Procedure Step Start Date Time
| 00404021 | Input Information Sequence
| 00404033 | Output Information Sequence
| 00404050 | Performed Procedure Step Start Date Time
| 00404051 | Performed Procedure Step End Date Time
| 00741000 | Procedure Step State | SCHEDULED<br/>IN PROGRESS<br/>CANCELED<br/>COMPLETED
| **00741204** | **Procedure Step Label** | **COVID-19**
| 00741216 | Unified Procedure Step Performed Procedure Sequence
| 00741238 | Reason For Cancellation | *Timeout*<sup>1</sup><br/>No Data<sup>2</sup><br/>Invalid Data<br/>Unknown Error

>1. Timeout: analysis doesn't finish in 10 minutes
>2. No Data: analysis request scheduled but no data received in 2 hours

### 2.3 Request Analysis

#### POST /workitems?{AffectedSOPInstanceUID}

- **API rate limit**: 5 times per minute

- **Response ETA**: 1 minute

##### AffectedSOPInstanceUID requirements

- **MUST** be composed with decimal numbers and dots
- **MUST** be encoded in UTF-8
- Length **MUST** be less than 64 characters long

#### Status codes

| Code | Description
| ---  | ---
| 201 Created | New analysis request has been created
| 400 Bad Request | [Required keys](#required-keys) missing or AffectedSOPInstanceUID doesn't meet the [requirements](#affectedsopinstanceuid-requirements)
| 401 Unauthorized | Token invalid or expired
| 409 Conflict | An analysis request with same ID already exists
| 503 Service Unavailable | The server is currently unable to handle the request due to a temporary overload

#### Required Keys

| Key | Description
| :---: | ---
| 00741204 | Procedure Step Label
| 00404021 | Input Information Sequence
| 0020000D | Study Instance UID

#### Request

```http
POST /workitems?1.3.45.214 HTTP/1.1
Content-Type: application/json
Authorization: Bearer {token}

{
  "00741204": "COVID-19",
  "00404021": {
    "0020000D": "2.16.840.1.113662.2.1.99999.5175439602988854"
  }
}
```

#### Response

```http
HTTP/1.1 201 Created
Location: /workitems?1.3.45.214
Content-Type: application/json
```

### 2.4 Query Status

#### GET /workitems/{AffectedSOPInstanceUID}

- **AffectedSOPInstanceUID requirements**: see [here](#affectedsopinstanceuid-requirements)

- **API rate limit**: 60 times per minute

- **Response ETA**: 1 second

#### Status codes

| Code | Description
| ---  | ---
| 200 OK | All good, please check result in response payload
| 401 Unauthorized | Token invalid or expired
| 404 Not Found | AffectedSOPInstanceUID provided doesn't exist
| 503 Service Unavailable | The server is currently unable to handle the request due to a temporary overload

#### Request

```http
GET /workitems/1.3.45.214 HTTP/1.1
Authorization: Bearer {token}
```

#### Response

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "00404005": {
    "Value": ["2020-03-26 05:21:35.67521964 +0000 UTC m=+8099.205393999"],
    "vr": "DT"
  },
  "00741000": {
    "Value": ["COMPLETED"],
    "vr": "CS"
  },
  "00741216": {
    "Value": [{
      "00400280": {
        "Value": ["Positive"],
        "vr": "ST"
      },
      "00404033": {
        "Value": [{
          "0020000D": {
            "Value": ["2.16.840.1.113662.2.1.1778.250300113532482162"],
            "vr": "UI"
          }
        }],
        "vr": "SQ"
      },
      "00404050": {
        "Value": ["2020-03-26 05:21:40.163982532 +0000 UTC m=+4981.844971167"],
        "vr": "DT"
      },
      "00404051": {
        "Value": ["2020-03-26 05:23:21.28836734 +0000 UTC m=+5082.969355969"],
        "vr": "DT"
      }
    }],
    "vr": "SQ"
  }
}
```

> See [here](#22-http-requestresponse-payload-syntax) for key explanation

#### Cancellation

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "00741000": {
    "Value": ["CANCELED"],
    "vr": "CS"
  },
  "00741238": {
    "Value": ["Timeout"],
    "vr": "LT"
  }
}
```

#### Completed

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "00741000": {
    "Value": ["COMPLETED"],
    "vr": "CS"
  },
  "00741216": {
    "Value": [{
      "00400280": {
        "Value": ["Positive"],
        "vr": "ST"
      }
    }],
    "vr": "SQ"
  }
}
```

## 3. DICOM Transfer

DICOM files **SHOULD** be sent to COVID-19 Server via DICOM `C-STORE` or `C-MOVE` command, for example

```shell
foo@bar:~$ dcmsend -v {ip} {port} -aec {target ae} -aet {source ae} +r +sd -nh /path/to/dicom
```

> dcmsend can be found in [DCMTK](https://support.dcmtk.org/docs/)

#### Requirments

- DICOM files **SHOULD** be de-identified before send to COVID-19 Server
- DICOM files **MUST** be transferred within 2 hours after analysis request triggered, otherwise the request will be canceled
- DICOM files will be cached for 24 hours by COVID-19 Detection Service

## 4. Free Trial Limitations

- The bearer token will be expired in 14 days
- DICOM communication bandwidth is 100Mbit/s
- AE title/ip/port of both sides **MUST** be exchanged and configured prior to the trial of this service

> The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be interpreted as described in [RFC 2119](https://tools.ietf.org/html/rfc2119)
