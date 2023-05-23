@api
Feature: Notification
  As a user
  I want to be notified of various events
  So that I can stay updated about the information

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Carol    |
    And user "Alice" has uploaded file with content "other data" to "/textfile1.txt"
    And user "Alice" has created folder "my_data"


  Scenario Outline: user gets a notification of resource sharing
    Given user "Alice" has shared entry "<resource>" with user "Brian"
    When user "Brian" lists all notifications
    Then the HTTP status code should be "200"
    And the JSON response should contain a notification message with the subject "Resource shared" and the message-details should match
    """
    {
    "type": "object",
      "required": [
        "app",
        "datetime",
        "message",
        "messageRich",
        "messageRichParameters",
        "notification_id",
        "object_id",
        "object_type",
        "subject",
        "subjectRich",
        "user"
      ],
      "properties": {
        "app": {
          "type": "string",
          "enum": ["userlog"]
        },
        "message": {
          "type": "string",
          "enum": ["Alice Hansen shared <resource> with you"]
        },
        "messageRich": {
          "type": "string",
          "enum": ["{user} shared {resource} with you"]
        },
        "messageRichParameters": {
          "type": "object",
          "required": [
            "resource",
            "user"
          ],
          "properties": {
            "resource": {
              "type": "object",
              "required": [
                "id",
                "name"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}\\$[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}![a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$"
                },
                "name": {
                  "type": "string",
                  "enum": ["<resource>"]
                }
              }
            },
            "user": {
              "type": "object",
              "required": [
                "displayname",
                "id",
                "name"
              ],
              "properties": {
                "displayname": {
                  "type": "string",
                  "enum": ["Alice Hansen"]
                },
                "id": {
                  "type": "string",
                  "enim": ["%user_id%"]
                },
                "name": {
                  "type": "string",
                  "enum": ["Alice"]
                }
              }
            }
          }
        },
        "notification_id": {
          "type": "string"
        },
        "object_id": {
          "type": "string"
        },
        "object_type": {
          "type": "string",
          "enum": ["share"]
        },
        "subject": {
          "type": "string",
          "enum": ["Resource shared"]
        },
        "subjectRich": {
          "type": "string",
          "enum": ["Resource shared"]
        },
        "user": {
          "type": "string",
          "enum": ["Alice"]
        }
      }
    }
    """
    Examples:
      | resource      |
      | textfile1.txt |
      | my_data       |


  Scenario Outline: user gets a notification of unsharing resource
    Given user "Alice" has shared entry "<resource>" with user "Brian"
    And user "Brian" has accepted share "/<resource>" offered by user "Alice"
    And user "Alice" has unshared entity "<resource>" shared to "Brian"
    When user "Brian" lists all notifications
    Then the HTTP status code should be "200"
    And the JSON response should contain a notification message with the subject "Resource unshared" and the message-details should match
    """
    {
    "type": "object",
      "required": [
        "app",
        "datetime",
        "message",
        "messageRich",
        "messageRichParameters",
        "notification_id",
        "object_id",
        "object_type",
        "subject",
        "subjectRich",
        "user"
      ],
      "properties": {
        "app": {
          "type": "string",
          "enum": ["userlog"]
        },
        "message": {
          "type": "string",
          "enum": ["Alice Hansen unshared <resource> with you"]
        },
        "messageRich": {
          "type": "string",
          "enum": ["{user} unshared {resource} with you"]
        },
        "messageRichParameters": {
          "type": "object",
          "required": [
            "resource",
            "user"
          ],
          "properties": {
            "resource": {
              "type": "object",
              "required": [
                "id",
                "name"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}\\$[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}![a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$"
                },
                "name": {
                  "type": "string",
                  "enum": ["<resource>"]
                }
              }
            },
            "user": {
              "type": "object",
              "required": [
                "displayname",
                "id",
                "name"
              ],
              "properties": {
                "displayname": {
                  "type": "string",
                  "enum": ["Alice Hansen"]
                },
                "id": {
                  "type": "string",
                  "enim": ["%user_id%"]
                },
                "name": {
                  "type": "string",
                  "enum": ["Alice"]
                }
              }
            }
          }
        },
        "notification_id": {
          "type": "string"
        },
        "object_id": {
          "type": "string"
        },
        "object_type": {
          "type": "string",
          "enum": ["share"]
        },
        "subject": {
          "type": "string",
          "enum": ["Resource unshared"]
        },
        "subjectRich": {
          "type": "string",
          "enum": ["Resource unshared"]
        },
        "user": {
          "type": "string",
          "enum": ["Alice"]
        }
      }
    }
    """
    Examples:
      | resource      |
      | textfile1.txt |
      | my_data       |
