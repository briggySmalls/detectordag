package iot

var provisioningTemplate = `
{
    "Parameters" : {
        "DeviceId": {
            "Type": "String"
        },
        "ThingGroup": {
            "Type": "String"
        },
        "ThingType": {
            "Type": "String"
        },
        "CertificateId": {
            "Type": "String"
        },
        "DeviceName": {
            "Type": "String"
        },
        "AccountId": {
            "Type": "String"
        }
    },
    "Resources": {
        "thing" : {
            "Type" : "AWS::IoT::Thing",
            "Properties" : {
                "ThingName": {"Ref": "DeviceId"},
                "ThingGroups": [{"Ref" : "ThingGroup"}],
                "ThingTypeName" :  {"Ref" : "ThingType"},
                "AttributePayload" : {"name": {"Ref": "DeviceName"}, "account-id": {"Ref": "AccountId"}, "visibility": "true"}
            }
        },
        "certificate" : {
            "Type" : "AWS::IoT::Certificate",
            "Properties" : {
                "CertificateId": {"Ref": "CertificateId"}
            }
        }
    }
}
`
