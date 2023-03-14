@api @skipOnOcV10
Feature: Change data of space
  As a user with space admin rights
  I want to be able to change the data of a created space (increase the quota, change name, etc.)

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Bob      |
    And the administrator has given "Alice" the role "Space Admin" using the settings api
    And user "Alice" has created a space "Project Jupiter" of type "project" with quota "20"
    And user "Alice" has shared a space "Project Jupiter" with settings:
      | shareWith | Brian    |
      | role      | editor |
    And user "Alice" has shared a space "Project Jupiter" with settings:
      | shareWith | Bob    |
      | role      | viewer |
    And using spaces DAV path


  Scenario Outline: Only space admin user can change the name of a Space via the Graph API
    When user "<user>" changes the name of the "Project Jupiter" space to "<expectedName>"
    Then the HTTP status code should be "<code>"
    And for user "<user>" the JSON response should contain space called "<expectedName>" and match
    """
     {
      "type": "object",
      "required": [
        "name"
      ],
      "properties": {
        "name": {
          "type": "string",
          "enum": ["<expectedName>"]
        }
      }
    }
    """
    Examples:
      | user  | code | expectedName       |
      | Alice | 200  | Project Death Star |
      | Brian | 403  | Project Jupiter    |
      | Bob   | 403  | Project Jupiter    |

  Scenario: Only space admin user can change the description(subtitle) of a Space via the Graph API
    When user "Alice" changes the description of the "Project Jupiter" space to "The Death Star is a fictional mobile space station"
    Then the HTTP status code should be "200"
    And for user "Alice" the JSON response should contain space called "Project Jupiter" and match
    """
     {
      "type": "object",
      "required": [
        "name",
        "driveType",
        "description"
      ],
      "properties": {
        "driveType": {
          "type": "string",
          "enum": ["project"]
        },
        "name": {
          "type": "string",
          "enum": ["Project Jupiter"]
        },
        "description": {
          "type": "string",
          "enum": ["The Death Star is a fictional mobile space station"]
        }
      }
    }
    """

  Scenario Outline: Viewer and editor cannot change the description(subtitle) of a Space via the Graph API
    When user "<user>" changes the description of the "Project Jupiter" space to "The Death Star is a fictional mobile space station"
    Then the HTTP status code should be "<code>"
    Examples:
      | user  | code |
      | Brian | 403  |
      | Bob   | 403  |


  Scenario Outline: An user tries to increase the quota of a Space via the Graph API
    When user "<user>" changes the quota of the "Project Jupiter" space to "100"
    Then the HTTP status code should be "<code>"
    And for user "<user>" the JSON response should contain space called "Project Jupiter" and match

    """
     {
      "type": "object",
      "required": [
        "name",
        "quota"
      ],
      "properties": {
        "name": {
          "type": "string",
          "enum": ["Project Jupiter"]
        },
        "quota": {
          "type": "object",
          "required": [
            "total"
          ],
          "properties": {
            "total" : {
              "type": "number",
              "enum": [<expectedQuataValue>]
            }
          }
        }
      }
    }
    """
    Examples:
      | user  | code | expectedQuataValue |
      | Alice | 200  | 100                |
      | Brian | 401  | 20                 |
      | Bob   | 401  | 20                 |


  Scenario Outline: An space admin user set no restriction quota of a Space via the Graph API
    When user "Alice" changes the quota of the "Project Jupiter" space to "<quotaValue>"
    Then the HTTP status code should be "200"
    When user "Alice" uploads a file inside space "Project Jupiter" with content "some content" to "file.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the JSON response should contain space called "Project Jupiter" and match
    """
     {
      "type": "object",
      "required": [
        "name",
        "quota"
      ],
      "properties": {
        "name": {
          "type": "string",
          "enum": ["Project Jupiter"]
        },
        "quota": {
          "type": "object",
          "required": [
            "used"
          ],
          "properties": {
            "used" : {
              "type": "number",
              "enum": [12]
            }
          }
        }
      }
    }
    """
    Examples:
      | quotaValue |
      | 0          |
      | -1         |


  Scenario: An user space admin set readme file as description of the space via the Graph API
    Given user "Alice" has created a folder ".space" in space "Project Jupiter"
    And user "Alice" has uploaded a file inside space "Project Jupiter" with content "space description" to ".space/readme.md"
    When user "Alice" sets the file ".space/readme.md" as a description in a special section of the "Project Jupiter" space
    Then the HTTP status code should be "200"
    And for user "Alice" the JSON response should contain space called "Project Jupiter" owned by "Alice" with description file ".space/readme.md" and match
    """
    {
      "type": "object",
      "required": [
        "name",
        "special"
      ],
      "properties": {
        "name": {
          "type": "string",
          "enum": ["Project Jupiter"]
        },
        "special": {
          "type": "array",
          "items": [
            {
              "type": "object",
              "required": [
                "size",
                "name",
                "specialFolder",
                "file",
                "id",
                "eTag"
              ],
              "properties": {
                "size": {
                  "type": "number",
                  "enum": [17]
                },
                "name": {
                  "type": "string",
                  "enum": ["readme.md"]
                },
                "specialFolder": {
                  "type": "object",
                  "required": [
                    "name"
                  ],
                  "properties": {
                    "name": {
                      "type": "string",
                      "enum": ["readme"]
                    }
                  }
                },
                "file": {
                  "type": "object",
                  "required": [
                    "mimeType"
                  ],
                  "properties": {
                    "type": "string",
                    "enum": ["text/markdown"]
                  }
                },
                "id": {
                  "type": "string",
                  "enum": ["%file_id%"]
                },
                "tag": {
                  "type": "string",
                  "enum": ["%eTag%"]
                }
              }
            }
          ]
        }
      }
    }
    """
    And for user "Alice" folder ".space/" of the space "Project Jupiter" should contain these entries:
      | readme.md |
    And for user "Alice" the content of the file ".space/readme.md" of the space "Project Jupiter" should be "space description"


  Scenario Outline: An user member of the space changes readme file
    Given user "Alice" has created a folder ".space" in space "Project Jupiter"
    And user "Alice" has uploaded a file inside space "Project Jupiter" with content "space description" to ".space/readme.md"
    And user "Alice" has set the file ".space/readme.md" as a description in a special section of the "Project Jupiter" space
    When user "<user>" uploads a file inside space "Project Jupiter" with content "new description" to ".space/readme.md" using the WebDAV API
    Then the HTTP status code should be "<code>"
    And for user "Alice" the JSON response should contain space called "Project Jupiter" owned by "Alice" with description file ".space/readme.md" and match
    """
    {
      "type": "object",
      "required": [
        "name",
        "special"
      ],
      "properties": {
        "name": {
          "type": "string",
          "enum": ["Project Jupiter"]
        },
        "special": {
          "type": "array",
          "items": [
            {
              "type": "object",
              "required": [
                "size",
                "name",
                "specialFolder",
                "file",
                "id",
                "eTag"
              ],
              "properties": {
                "size": {
                  "type": "number",
                  "enum": [<size>]
                },
                "name": {
                  "type": "string",
                  "enum": ["readme.md"]
                },
                "specialFolder": {
                  "type": "object",
                  "required": [
                    "name"
                  ],
                  "properties": {
                    "name": {
                      "type": "string",
                      "enum": ["readme"]
                    }
                  }
                },
                "file": {
                  "type": "object",
                  "required": [
                    "mimeType"
                  ],
                  "properties": {
                    "type": "string",
                    "enum": ["text/markdown"]
                  }
                },
                "id": {
                  "type": "string",
                  "enum": ["%file_id%"]
                },
                "tag": {
                  "type": "string",
                  "enum": ["%eTag%"]
                }
              }
            }
          ]
        }
      }
    }
    """
    And for user "<user>" folder ".space/" of the space "Project Jupiter" should contain these entries:
      | readme.md |
    And for user "<user>" the content of the file ".space/readme.md" of the space "Project Jupiter" should be "<content>"
    Examples:
      | user  | code | size | content           |
      | Brian | 204  | 15   | new description   |
      | Bob   | 403  | 17   | space description |


  Scenario Outline: An user space admin and editor set image file as space image of the space via the Graph API
    Given user "Alice" has created a folder ".space" in space "Project Jupiter"
    And user "<user>" has uploaded a file inside space "Project Jupiter" with content "" to ".space/<fileName>"
    When user "<user>" sets the file ".space/<fileName>" as a space image in a special section of the "Project Jupiter" space
    Then the HTTP status code should be "200"
    And for user "Alice" the JSON response should contain space called "Project Jupiter" owned by "Alice" with description file ".space/<fileName>" and match
    """
    {
      "type": "object",
      "required": [
        "name",
        "special"
      ],
      "properties": {
        "name": {
          "type": "string",
          "enum": ["Project Jupiter"]
        },
        "special": {
          "type": "array",
          "items": [
            {
              "type": "object",
              "required": [
                "size",
                "name",
                "specialFolder",
                "file",
                "id",
                "eTag"
              ],
              "properties": {
                "size": {
                  "type": "number",
                  "enum": [0]
                },
                "name": {
                  "type": "string",
                  "enum": ["<nameInResponse>"]
                },
                "specialFolder": {
                  "type": "object",
                  "required": [
                    "name"
                  ],
                  "properties": {
                    "name": {
                      "type": "string",
                      "enum": ["image"]
                    }
                  }
                },
                "file": {
                  "type": "object",
                  "required": [
                    "mimeType"
                  ],
                  "properties": {
                    "type": "string",
                    "enum": ["<mimeType>"]
                  }
                },
                "id": {
                  "type": "string",
                  "enum": ["%file_id%"]
                },
                "tag": {
                  "type": "string",
                  "enum": ["%eTag%"]
                }
              }
            }
          ]
        }
      }
    }
    """
    And for user "<user>" folder ".space/" of the space "Project Jupiter" should contain these entries:
      | <fileName> |
    Examples:
      | user  | fileName        | nameInResponse  | mimeType   |
      | Alice | spaceImage.jpeg | spaceImage.jpeg | image/jpeg |
      | Brian | spaceImage.png  | spaceImage.png  | image/png  |
      | Alice | spaceImage.gif  | spaceImage.gif  | image/gif  |

  Scenario: An user viewer cannot set image file as space image of the space via the Graph API
    Given user "Alice" has created a folder ".space" in space "Project Jupiter"
    And user "Alice" has uploaded a file inside space "Project Jupiter" with content "" to ".space/someImageFile.jpg"
    When user "Bob" sets the file ".space/someImageFile.jpg" as a space image in a special section of the "Project Jupiter" space
    Then the HTTP status code should be "403"


  Scenario Outline: An user set new readme file as description of the space via the Graph API
    Given user "Alice" has created a folder ".space" in space "Project Jupiter"
    And user "Alice" has uploaded a file inside space "Project Jupiter" with content "space description" to ".space/readme.md"
    And user "Alice" has set the file ".space/readme.md" as a description in a special section of the "Project Jupiter" space
    When user "<user>" uploads a file inside space "Project Jupiter" owned by the user "Alice" with content "new content" to ".space/readme.md" using the WebDAV API
    Then the HTTP status code should be "<code>"
    And for user "<user>" the content of the file ".space/readme.md" of the space "Project Jupiter" should be "<expectedContent>"
    And for user "<user>" the JSON response should contain space called "Project Jupiter" owned by "Alice" with description file ".space/readme.md" and match
    """
    {
      "type": "object",
      "required": [
        "name",
        "special"
      ],
      "properties": {
        "name": {
          "type": "string",
          "enum": ["Project Jupiter"]
        },
        "special": {
          "type": "array",
          "items": [
            {
              "type": "object",
              "required": [
                "size",
                "name",
                "specialFolder",
                "file",
                "id",
                "eTag"
              ],
              "properties": {
                "size": {
                  "type": "number",
                  "enum": [<expectedSize>]
                },
                "name": {
                  "type": "string",
                  "enum": ["readme.md"]
                },
                "specialFolder": {
                  "type": "object",
                  "required": [
                    "name"
                  ],
                  "properties": {
                    "name": {
                      "type": "string",
                      "enum": ["readme"]
                    }
                  }
                },
                "file": {
                  "type": "object",
                  "required": [
                    "mimeType"
                  ],
                  "properties": {
                    "type": "string",
                    "enum": ["text/markdown"]
                  }
                },
                "id": {
                  "type": "string",
                  "enum": ["%file_id%"]
                },
                "tag": {
                  "type": "string",
                  "enum": ["%eTag%"]
                }
              }
            }
          ]
        }
      }
    }
    """
    Examples:
      | user  | code | expectedSize | expectedContent   |
      | Alice | 204  | 11           | new content       |
      | Brian | 204  | 11           | new content       |
      | Bob   | 403  | 17           | space description |


  Scenario Outline: An user set new image file as space image of the space via the Graph API
    Given user "Alice" has created a folder ".space" in space "Project Jupiter"
    And user "Alice" has uploaded a file inside space "Project Jupiter" with content "" to ".space/spaceImage.jpeg"
    And user "Alice" has set the file ".space/spaceImage.jpeg" as a space image in a special section of the "Project Jupiter" space
    When user "<user>" has uploaded a file inside space "Project Jupiter" with content "" to ".space/newSpaceImage.png"
    And user "<user>" sets the file ".space/newSpaceImage.png" as a space image in a special section of the "Project Jupiter" space
    Then the HTTP status code should be "200"
    And for user "<user>" the JSON response should contain space called "Project Jupiter" owned by "Alice" with description file ".space/newSpaceImage.png" and match
    """
    {
      "type": "object",
      "required": [
        "name",
        "special"
      ],
      "properties": {
        "name": {
          "type": "string",
          "enum": ["Project Jupiter"]
        },
        "special": {
          "type": "array",
          "items": [
            {
              "type": "object",
              "required": [
                "size",
                "name",
                "specialFolder",
                "file",
                "id",
                "eTag"
              ],
              "properties": {
                "size": {
                  "type": "number",
                  "enum": [0]
                },
                "name": {
                  "type": "string",
                  "enum": ["newSpaceImage.png"]
                },
                "specialFolder": {
                  "type": "object",
                  "required": [
                    "name"
                  ],
                  "properties": {
                    "name": {
                      "type": "string",
                      "enum": ["image"]
                    }
                  }
                },
                "file": {
                  "type": "object",
                  "required": [
                    "mimeType"
                  ],
                  "properties": {
                    "type": "string",
                    "enum": ["image/png"]
                  }
                },
                "id": {
                  "type": "string",
                  "enum": ["%file_id%"]
                },
                "tag": {
                  "type": "string",
                  "enum": ["%eTag%"]
                }
              }
            }
          ]
        }
      }
    }
    """
    Examples:
      | user  |
      | Alice |
      | Brian |


  Scenario Outline: An admin user set own quota of a personal space via the Graph API
    When user "Admin" changes the quota of the "Admin" space to "<quotaValue>"
    Then the HTTP status code should be "200"
    When user "Admin" uploads a file inside space "Admin" with content "file is more than 15 bytes" to "file.txt" using the WebDAV API
    Then the HTTP status code should be <code>
    Examples:
      | quotaValue | code                    |
      | 15         | "507"                   |
      | 10000      | between "201" and "204" |
      | 0          | between "201" and "204" |
      | -1         | between "201" and "204" |


  Scenario Outline: An admin user set an user personal space quota of via the Graph API
    When user "Admin" changes the quota of the "Brian Murphy" space to "<quotaValue>"
    Then the HTTP status code should be "200"
    When user "Brian" uploads a file inside space "Brian Murphy" with content "file is more than 15 bytes" to "file.txt" using the WebDAV API
    Then the HTTP status code should be <code>
    Then the user "Brian" should have a space called "Brian Murphy" with these key and value pairs:
      | key           | value   |
      | quota@@@total | <total> |
      | quota@@@used  | <used>  |
    Examples:
      | quotaValue | code                    | total | used |
      | 15         | "507"                   | 15    | 0    |
      | 10000      | between "201" and "204" | 10000 | 26   |
      | 0          | between "201" and "204" | 0     | 26   |
      | -1         | between "201" and "204" | 0     | 26   |


  Scenario: user sends invalid space uuid via the Graph API
    When user "Admin" tries to change the name of the "non-existing" space to "new name"
    Then the HTTP status code should be "404"
    When user "Admin" tries to change the quota of the "non-existing" space to "10"
    Then the HTTP status code should be "404"
    When user "Alice" tries to change the description of the "non-existing" space to "new description"
    Then the HTTP status code should be "404"


  Scenario: user sends PATCH request to other user's space that they don't have access to
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Carol    |
    When user "Carol" sends PATCH request to the space "Personal" of user "Alice" with data "{}"
    Then the HTTP status code should be "404"
    When user "Carol" sends PATCH request to the space "Project Jupiter" of user "Alice" with data "{}"
    Then the HTTP status code should be "404"
