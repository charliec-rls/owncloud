Feature: Send a sharing invitations
  As the owner of a resource
  I want to be able to send invitations to other users
  So that they can have access to it

  https://owncloud.dev/libre-graph-api/#/drives.permissions/Invite

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |


  Scenario Outline: send share invitation to user with different roles
    Given user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following share invitation using the Graph API:
      | resource        | <resource>         |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "200"
    And for user "Brian" the space Shares should contain these entries:
      | <resource> |
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "value"
        ],
        "properties": {
          "value": {
            "type": "array",
            "maxItems": 1,
            "minItems": 1,
            "items": {
              "type": "object",
              "required": [
                "id",
                "roles",
                "grantedToV2"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "maxItems": 1,
                  "minItems": 1,
                  "items": {
                    "type": "string",
                    "pattern": "^%role_id_pattern%$"
                  }
                },
                "grantedToV2": {
                  "type": "object",
                  "required": [
                    "user"
                  ],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": [
                        "id",
                        "displayName"
                      ],
                      "properties": {
                        "id": {
                          "type": "string",
                          "pattern": "^%user_id_pattern%$"
                        },
                        "displayName": {
                          "type": "string",
                          "enum": [
                            "Brian Murphy"
                          ]
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource       |
      | Viewer           | /textfile1.txt |
      | File Editor      | /textfile1.txt |
      | Viewer           | FolderToShare  |
      | Editor           | FolderToShare  |
      | Uploader         | FolderToShare  |


  Scenario Outline: send share invitation to group with different roles
    Given user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following share invitation using the Graph API:
      | resource        | <resource>         |
      | space           | Personal           |
      | sharee          | grp1               |
      | shareType       | group              |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "200"
    And for user "Brian" the space Shares should contain these entries:
      | <resource> |
    And for user "Carol" the space Shares should contain these entries:
      | <resource> |
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "value"
        ],
        "properties": {
          "value": {
            "type": "array",
            "maxItems": 1,
            "minItems": 1,
            "items": {
              "type": "object",
              "required": [
                "id",
                "roles",
                "grantedToV2"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "maxItems": 1,
                  "minItems": 1,
                  "items": {
                    "type": "string",
                    "pattern": "^%role_id_pattern%$"
                  }
                },
                "grantedToV2": {
                  "type": "object",
                  "required": [
                    "group"
                  ],
                  "properties": {
                    "group": {
                      "type": "object",
                      "required": [
                        "id",
                        "displayName"
                      ],
                      "properties": {
                        "id": {
                          "type": "string",
                          "pattern": "^%group_id_pattern%$"
                        },
                        "displayName": {
                          "type": "string",
                          "enum": [
                            "grp1"
                          ]
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource       |
      | Viewer           | /textfile1.txt |
      | File Editor      | /textfile1.txt |
      | Viewer           | FolderToShare  |
      | Editor           | FolderToShare  |
      | Uploader         | FolderToShare  |


  Scenario Outline: send share invitation for a file to user with different permissions
    Given user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    When user "Alice" sends the following share invitation using the Graph API:
      | resource          | textfile1.txt        |
      | space             | Personal             |
      | sharee            | Brian                |
      | shareType         | user                 |
      | permissionsAction | <permissions-action> |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "value"
        ],
        "properties": {
          "value": {
            "type": "array",
            "maxItems": 1,
            "minItems": 1,
            "items": {
              "type": "object",
              "required": [
                "id",
                "@libre.graph.permissions.actions",
                "grantedToV2"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "@libre.graph.permissions.actions": {
                  "type": "array",
                  "maxItems": 1,
                  "minItems": 1,
                  "items": {
                    "type": "string",
                    "pattern": "^libre\\.graph\\/driveItem\\/<permissions-action>$"
                  }
                },
                "grantedToV2": {
                  "type": "object",
                  "required": [
                    "user"
                  ],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": [
                        "id",
                        "displayName"
                      ],
                      "properties": {
                        "id": {
                          "type": "string",
                          "pattern": "^%user_id_pattern%$"
                        },
                        "displayName": {
                          "type": "string",
                          "enum": [
                            "Brian Murphy"
                          ]
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-action |
      | upload/create      |
      | path/read          |
      | quota/read         |
      | content/read       |
      | permissions/read   |
      | children/read      |
      | versions/read      |
      | deleted/read       |
      | basic/read         |
      | versions/update    |
      | deleted/update     |
      | deleted/delete     |


  Scenario Outline: send share invitation for a folder to user with different permissions
    Given user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following share invitation using the Graph API:
      | resource          | FolderToShare        |
      | space             | Personal             |
      | sharee            | Brian                |
      | shareType         | user                 |
      | permissionsAction | <permissions-action> |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "value"
        ],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": [
                "id",
                "@libre.graph.permissions.actions",
                "grantedToV2"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "@libre.graph.permissions.actions": {
                  "type": "array",
                  "minItems": 1,
                  "maxItems": 1,
                  "items": {
                    "type": "string",
                    "pattern": "^libre\\.graph\\/driveItem\\/<permissions-action>$"
                  }
                },
                "grantedToV2": {
                  "type": "object",
                  "required": [
                    "user"
                  ],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": [
                        "id",
                        "displayName"
                      ],
                      "properties": {
                        "id": {
                          "type": "string",
                          "pattern": "^%user_id_pattern%$"
                        },
                        "displayName": {
                          "type": "string",
                          "enum": [
                            "Brian Murphy"
                          ]
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-action |
      | children/create    |
      | upload/create      |
      | path/read          |
      | quota/read         |
      | content/read       |
      | permissions/read   |
      | children/read      |
      | versions/read      |
      | deleted/read       |
      | basic/read         |
      | path/update        |
      | versions/update    |
      | deleted/update     |
      | standard/delete    |
      | deleted/delete     |


  Scenario Outline: send share invitation for a file to group with different permissions
    Given user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    When user "Alice" sends the following share invitation using the Graph API:
      | resource          | textfile1.txt        |
      | space             | Personal             |
      | sharee            | grp1                 |
      | shareType         | group                |
      | permissionsAction | <permissions-action> |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "value"
        ],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": [
                "id",
                "@libre.graph.permissions.actions",
                "grantedToV2"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "@libre.graph.permissions.actions": {
                  "type": "array",
                  "minItems": 1,
                  "maxItems": 1,
                  "items": {
                    "type": "string",
                    "pattern": "^libre\\.graph\\/driveItem\\/<permissions-action>$"
                  }
                },
                "grantedToV2": {
                  "type": "object",
                  "required": [
                    "group"
                  ],
                  "properties": {
                    "group": {
                      "type": "object",
                      "required": [
                        "id",
                        "displayName"
                      ],
                      "properties": {
                        "id": {
                          "type": "string",
                          "pattern": "^%user_id_pattern%$"
                        },
                        "displayName": {
                          "type": "string",
                          "enum": [
                            "grp1"
                          ]
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-action |
      | upload/create      |
      | path/read          |
      | quota/read         |
      | content/read       |
      | permissions/read   |
      | children/read      |
      | versions/read      |
      | deleted/read       |
      | basic/read         |
      | versions/update    |
      | deleted/update     |
      | deleted/delete     |


  Scenario Outline: send share invitation for a folder to group with different permissions
    Given user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following share invitation using the Graph API:
      | resource          | FolderToShare        |
      | space             | Personal             |
      | sharee            | grp1                 |
      | shareType         | group                |
      | permissionsAction | <permissions-action> |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "value"
        ],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": [
                "id",
                "@libre.graph.permissions.actions",
                "grantedToV2"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "@libre.graph.permissions.actions": {
                  "type": "array",
                  "minItems": 1,
                  "maxItems": 1,
                  "items": {
                    "type": "string",
                    "pattern": "^libre\\.graph\\/driveItem\\/<permissions-action>$"
                  }
                },
                "grantedToV2": {
                  "type": "object",
                  "required": [
                    "group"
                  ],
                  "properties": {
                    "group": {
                      "type": "object",
                      "required": [
                        "id",
                        "displayName"
                      ],
                      "properties": {
                        "id": {
                          "type": "string",
                          "pattern": "^%user_id_pattern%$"
                        },
                        "displayName": {
                          "type": "string",
                          "enum": [
                            "grp1"
                          ]
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-action |
      | children/create    |
      | upload/create      |
      | path/read          |
      | quota/read         |
      | content/read       |
      | permissions/read   |
      | children/read      |
      | versions/read      |
      | deleted/read       |
      | basic/read         |
      | path/update        |
      | versions/update    |
      | deleted/update     |
      | standard/delete    |
      | deleted/delete     |


  Scenario Outline: send share invitation with expiration date to user with different roles
    Given user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following share invitation using the Graph API:
      | resource        | <resource>               |
      | space           | Personal                 |
      | sharee          | Brian                    |
      | shareType       | user                     |
      | permissionsRole | <permissions-role>       |
      | expireDate      | 2043-07-15T14:00:00.000Z |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "value"
        ],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": [
                "id",
                "roles",
                "grantedToV2",
                "expirationDateTime"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "minItems": 1,
                  "maxItems": 1,
                  "items": {
                    "type": "string",
                    "pattern": "^%role_id_pattern%$"
                  }
                },
                "grantedToV2": {
                  "type": "object",
                  "required": [
                    "user"
                  ],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": [
                        "id",
                        "displayName"
                      ],
                      "properties": {
                        "id": {
                          "type": "string",
                          "pattern": "^%user_id_pattern%$"
                        },
                        "displayName": {
                          "type": "string",
                          "enum": [
                            "Brian Murphy"
                          ]
                        }
                      }
                    }
                  }
                },
                "expirationDateTime": {
                  "type": "string",
                  "enum": [
                    "2043-07-15T14:00:00Z"
                  ]
                }
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource       |
      | Viewer           | /textfile1.txt |
      | File Editor      | /textfile1.txt |
      | Viewer           | FolderToShare  |
      | Editor           | FolderToShare  |
      | Uploader         | FolderToShare  |


  Scenario Outline: send share invitation with expiration date to group with different roles
    Given user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following share invitation using the Graph API:
      | resource        | <resource>               |
      | space           | Personal                 |
      | sharee          | grp1                     |
      | shareType       | group                    |
      | permissionsRole | <permissions-role>       |
      | expireDate      | 2043-07-15T14:00:00.000Z |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "value"
        ],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": [
                "id",
                "roles",
                "grantedToV2",
                "expirationDateTime"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "minItems": 1,
                  "maxItems": 1,
                  "items": {
                    "type": "string",
                    "pattern": "^%role_id_pattern%$"
                  }
                },
                "grantedToV2": {
                  "type": "object",
                  "required": [
                    "group"
                  ],
                  "properties": {
                    "group": {
                      "type": "object",
                      "required": [
                        "id",
                        "displayName"
                      ],
                      "properties": {
                        "id": {
                          "type": "string",
                          "pattern": "^%group_id_pattern%$"
                        },
                        "displayName": {
                          "type": "string",
                          "enum": [
                            "grp1"
                          ]
                        }
                      }
                    }
                  }
                },
                "expirationDateTime": {
                  "type": "string",
                  "enum": [
                    "2043-07-15T14:00:00Z"
                  ]
                }
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource       |
      | Viewer           | /textfile1.txt |
      | File Editor      | /textfile1.txt |
      | Viewer           | FolderToShare  |
      | Editor           | FolderToShare  |
      | Uploader         | FolderToShare  |

  @issue-7962
  Scenario Outline: send share invitation to disabled user
    Given user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    And the user "Admin" has disabled user "Brian"
    When user "Alice" sends the following share invitation using the Graph API:
      | resource        | <resource>         |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "value"
        ],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": [
                "id",
                "roles",
                "grantedToV2"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "minItems": 1,
                  "maxItems": 1,
                  "items": {
                    "type": "string",
                    "pattern": "^%role_id_pattern%$"
                  }
                },
                "grantedToV2": {
                  "type": "object",
                  "required": [
                    "user"
                  ],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": [
                        "id",
                        "displayName"
                      ],
                      "properties": {
                        "id": {
                          "type": "string",
                          "pattern": "^%user_id_pattern%$"
                        },
                        "displayName": {
                          "type": "string",
                          "enum": [
                            "Brian Murphy"
                          ]
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource       |
      | Viewer           | /textfile1.txt |
      | File Editor      | /textfile1.txt |
      | Viewer           | FolderToShare  |
      | Editor           | FolderToShare  |
      | Uploader         | FolderToShare  |


  Scenario Outline: send sharing invitation to a deleted group with different roles
    Given user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    And group "grp1" has been deleted
    When user "Alice" sends the following share invitation using the Graph API:
      | resource        | <resource>         |
      | space           | Personal           |
      | sharee          | grp1               |
      | shareType       | group              |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "error"
        ],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "message"
            ],
            "properties": {
              "code": {
                "type": "string",
                "enum": [
                  "invalidRequest"
                ]
              },
              "message": {
                "type": "string",
                "enum": [
                  "itemNotFound: not found"
                ]
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource       |
      | Viewer           | /textfile1.txt |
      | File Editor      | /textfile1.txt |
      | Viewer           | FolderToShare  |
      | Editor           | FolderToShare  |
      | Uploader         | FolderToShare  |


  Scenario Outline: send share invitation to deleted user
    Given user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    And the user "Admin" has deleted a user "Brian"
    When user "Alice" sends the following share invitation using the Graph API:
      | resource        | <resource>         |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "error"
        ],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "message"
            ],
            "properties": {
              "code": {
                "type": "string",
                "pattern": "invalidRequest"
              },
              "message": {
                "type": "string",
                "enum": [
                  "itemNotFound: not found"
                ]
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource       |
      | Viewer           | /textfile1.txt |
      | File Editor      | /textfile1.txt |
      | Viewer           | FolderToShare  |
      | Editor           | FolderToShare  |
      | Uploader         | FolderToShare  |


  Scenario Outline: try to send sharing invitation to multiple groups
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Carol    |
      | Bob      |
    And group "grp1" has been created
    And group "grp2" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp2      |
      | Bob      | grp2      |
    And user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following share invitation using the Graph API:
      | resource        | <resource>         |
      | space           | Personal           |
      | sharee          | grp1, grp2         |
      | shareType       | group, group       |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "error"
        ],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "message"
            ],
            "properties": {
              "code": {
                "type": "string",
                "enum": [
                  "invalidRequest"
                ]
              },
              "message": {
                "type": "string",
                "enum": [
                  "Key: 'DriveItemInvite.Recipients' Error:Field validation for 'Recipients' failed on the 'len' tag"
                ]
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource       |
      | Viewer           | /textfile1.txt |
      | File Editor      | /textfile1.txt |
      | Viewer           | FolderToShare  |
      | Editor           | FolderToShare  |
      | Uploader         | FolderToShare  |
      | Manager          | FolderToShare  |


  Scenario Outline: try to send sharing invitation to user and group at once
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Carol    |
      | Bob      |
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following share invitation using the Graph API:
      | resource        | <resource>         |
      | space           | Personal           |
      | sharee          | grp1, Bob          |
      | shareType       | group, user        |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "error"
        ],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "message"
            ],
            "properties": {
              "code": {
                "type": "string",
                "enum": [
                  "invalidRequest"
                ]
              },
              "message": {
                "type": "string",
                "enum": [
                  "Key: 'DriveItemInvite.Recipients' Error:Field validation for 'Recipients' failed on the 'len' tag"
                ]
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource       |
      | Viewer           | /textfile1.txt |
      | File Editor      | /textfile1.txt |
      | Viewer           | FolderToShare  |
      | Editor           | FolderToShare  |
      | Uploader         | FolderToShare  |


  Scenario Outline: send sharing invitation to non-existing group
    Given user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following share invitation using the Graph API:
      | resource        | <resource>         |
      | space           | Personal           |
      | sharee          | nonExistentGroup   |
      | shareType       | group              |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "error"
        ],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "message"
            ],
            "properties": {
              "code": {
                "type": "string",
                "enum": [
                  "invalidRequest"
                ]
              },
              "message": {
                "type": "string",
                "enum": [
                  "itemNotFound: not found"
                ]
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource       |
      | Viewer           | /textfile1.txt |
      | File Editor      | /textfile1.txt |
      | Viewer           | FolderToShare  |
      | Editor           | FolderToShare  |
      | Uploader         | FolderToShare  |


  Scenario Outline: send sharing invitation to already shared group
    Given user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    And user "Alice" has sent the following share invitation:
      | resource        | <resource>         |
      | space           | Personal           |
      | sharee          | grp1               |
      | shareType       | group              |
      | permissionsRole | <permissions-role> |
    When user "Alice" sends the following share invitation using the Graph API:
      | resource        | <resource>         |
      | space           | Personal           |
      | sharee          | grp1               |
      | shareType       | group              |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "409"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "error"
        ],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "message"
            ],
            "properties": {
              "code": {
                "type": "string",
                "enum": [
                  "nameAlreadyExists"
                ]
              },
              "message": {
                "type": "string",
                "pattern": "^error creating share: error: already exists:.*$"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource       |
      | Viewer           | /textfile1.txt |
      | File Editor      | /textfile1.txt |
      | Viewer           | FolderToShare  |
      | Editor           | FolderToShare  |
      | Uploader         | FolderToShare  |


  Scenario Outline: send share invitation to wrong user id
    Given user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" tries to send the following share invitation using the Graph API:
      | resource        | <resource>                           |
      | space           | Personal                             |
      | shareeId        | a4c0c83e-ae24-4870-93c3-fcaf2a2228f7 |
      | shareType       | user                                 |
      | permissionsRole | Viewer                               |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "error"
        ],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "message"
            ],
            "properties": {
              "code": {
                "type": "string",
                "enum": ["invalidRequest"]
              },
              "message": {
                "type": "string",
                "enum": [
                  "itemNotFound: not found"
                ]
              }
            }
          }
        }
      }
      """
    Examples:
      | resource       |
      | /textfile1.txt |
      | FolderToShare  |


  Scenario Outline: send share invitation with empty user id
    Given user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" tries to send the following share invitation using the Graph API:
      | resource        | <resource> |
      | space           | Personal   |
      | shareeId        |            |
      | shareType       | user       |
      | permissionsRole | Viewer     |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "error"
        ],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "message"
            ],
            "properties": {
              "code": {
                "type": "string",
                "enum": ["invalidRequest"]
              },
              "message": {
                "type": "string",
                "enum": [
                  "Key: 'DriveItemInvite.Recipients[0].ObjectId' Error:Field validation for 'ObjectId' failed on the 'ne' tag"
                ]
              }
            }
          }
        }
      }
      """
    Examples:
      | resource       |
      | /textfile1.txt |
      | FolderToShare  |


  Scenario Outline: send share invitation to user with wrong recipient type
    Given user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" tries to send the following share invitation using the Graph API:
      | resource        | <resource>     |
      | space           | Personal       |
      | sharee          | Brian          |
      | shareType       | wrongShareType |
      | permissionsRole | Viewer         |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "error"
        ],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "message"
            ],
            "properties": {
              "code": {
                "type": "string",
                "enum": ["invalidRequest"]
              },
              "message": {
                "type": "string",
                "enum": [
                  "Key: 'DriveItemInvite.Recipients[0].LibreGraphRecipientType' Error:Field validation for 'LibreGraphRecipientType' failed on the 'oneof' tag"
                ]
              }
            }
          }
        }
      }
      """
    Examples:
      | resource       |
      | /textfile1.txt |
      | FolderToShare  |


  Scenario Outline: send share invitation to group with wrong recipient type
    Given user "Carol" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    When user "Alice" tries to send the following share invitation using the Graph API:
      | resource        | <resource>     |
      | space           | Personal       |
      | sharee          | grp1           |
      | shareType       | wrongShareType |
      | permissionsRole | Viewer         |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "error"
        ],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "message"
            ],
            "properties": {
              "code": {
                "type": "string",
                "enum": ["invalidRequest"]
              },
              "message": {
                "type": "string",
                "enum": [
                  "Key: 'DriveItemInvite.Recipients[0].LibreGraphRecipientType' Error:Field validation for 'LibreGraphRecipientType' failed on the 'oneof' tag"
                ]
              }
            }
          }
        }
      }
      """
    Examples:
      | resource       |
      | /textfile1.txt |
      | FolderToShare  |


  Scenario Outline: send share invitation to user with empty recipient type
    Given user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" tries to send the following share invitation using the Graph API:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | Brian      |
      | shareType       |            |
      | permissionsRole | Viewer     |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "error"
        ],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "message"
            ],
            "properties": {
              "code": {
                "type": "string",
                "enum": ["invalidRequest"]
              },
              "message": {
                "type": "string",
                "enum": [
                  "Key: 'DriveItemInvite.Recipients[0].LibreGraphRecipientType' Error:Field validation for 'LibreGraphRecipientType' failed on the 'oneof' tag"
                ]
              }
            }
          }
        }
      }
      """
    Examples:
      | resource       |
      | /textfile1.txt |
      | FolderToShare  |


  Scenario Outline: send share invitation to group with empty recipient type
    Given user "Carol" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    When user "Alice" tries to send the following share invitation using the Graph API:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | grp1       |
      | shareType       |            |
      | permissionsRole | Viewer     |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "error"
        ],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "message"
            ],
            "properties": {
              "code": {
                "type": "string",
                "enum": ["invalidRequest"]
              },
              "message": {
                "type": "string",
                "enum": [
                  "Key: 'DriveItemInvite.Recipients[0].LibreGraphRecipientType' Error:Field validation for 'LibreGraphRecipientType' failed on the 'oneof' tag"
                ]
              }
            }
          }
        }
      }
      """
    Examples:
      | resource       |
      | /textfile1.txt |
      | FolderToShare  |


  Scenario Outline: try to share a resource with invalid roles
    Given user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following share invitation using the Graph API:
      | resource        | <resource>         |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "error"
        ],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "message"
            ],
            "properties": {
              "code": {
                "type": "string",
                "enum": [
                  "invalidRequest"
                ]
              },
              "message": {
                "type": "string",
                "enum": [
                  "role not applicable to this resource"
                ]
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource       |
      | Manager          | /textfile1.txt |
      | Space Viewer     | /textfile1.txt |
      | Space Editor     | /textfile1.txt |
      | Manager          | FolderToShare  |
      | Space Viewer     | FolderToShare  |
      | Space Editor     | FolderToShare  |


  Scenario Outline: try to share a file with invalid roles
    Given user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    When user "Alice" sends the following share invitation using the Graph API:
      | resource        | textfile1.txt      |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "error"
        ],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "message"
            ],
            "properties": {
              "code": {
                "type": "string",
                "enum": [
                  "invalidRequest"
                ]
              },
              "message": {
                "type": "string",
                "enum": [
                  "role not applicable to this resource"
                ]
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role |
      | Editor           |
      | Uploader         |


  Scenario Outline: send share invitation to already shared user
    Given user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    And user "Alice" has sent the following share invitation:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | Brian      |
      | shareType       | user       |
      | permissionsRole | Viewer     |
    When user "Alice" tries to send the following share invitation using the Graph API:
      | resource        | <resource> |
      | space           | Personal   |
      | sharee          | Brian      |
      | shareType       | user       |
      | permissionsRole | Viewer     |
    Then the HTTP status code should be "409"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "error"
        ],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "message"
            ],
            "properties": {
              "code": {
                "type": "string",
                "enum": ["nameAlreadyExists"]
              },
              "message": {
                "type": "string",
                "pattern": "^error creating share: error: already exists:.*$"
              }
            }
          }
        }
      }
      """
    Examples:
      | resource       |
      | /textfile1.txt |
      | FolderToShare  |


  Scenario Outline: send share invitation for project space to user with different roles
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    When user "Alice" sends the following share invitation for space using the Graph API:
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "value"
        ],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": [
                "grantedToV2",
                "roles"
              ],
              "properties": {
                "grantedToV2": {
                  "type": "object",
                  "required": [
                    "user"
                  ],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": [
                        "displayName",
                        "id"
                      ],
                      "properties": {
                        "displayName": {
                          "type": "string",
                          "enum": [
                            "Brian Murphy"
                          ]
                        },
                        "id": {
                          "type": "string",
                          "pattern": "^%user_id_pattern%$"
                        }
                      }
                    }
                  }
                },
                "roles": {
                  "type": "array",
                  "minItems": 1,
                  "maxItems": 1,
                  "items": {
                    "type": "string",
                    "pattern": "^%role_id_pattern%$"
                  }
                }
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |


  Scenario Outline: send share invitation for disabled project space to user with different roles
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Admin" has disabled a space "NewSpace"
    When user "Alice" sends the following share invitation for space using the Graph API:
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "404"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "error"
        ],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "message"
            ],
            "properties": {
              "code": {
                "type": "string",
                "enum": ["itemNotFound"]
              },
              "message": {
                "type": "string",
                "pattern": "^stat: error: not found: %user_id_pattern%$"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |


  Scenario Outline: send share invitation for deleted project space to user with different roles
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Admin" has disabled a space "NewSpace"
    And user "Admin" has deleted a space "NewSpace"
    When user "Alice" sends the following share invitation for space using the Graph API:
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "404"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "error"
        ],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "message"
            ],
            "properties": {
              "code": {
                "type": "string",
                "enum": ["itemNotFound"]
              },
              "message": {
                "type": "string",
                "enum": ["stat: error: not found: "]
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |


  Scenario Outline: send share invitation for project space to group with different roles
    Given using spaces DAV path
    And user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    When user "Alice" sends the following share invitation for space using the Graph API:
      | space           | NewSpace           |
      | sharee          | grp1               |
      | shareType       | group              |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "value"
        ],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": [
                "grantedToV2",
                "roles"
              ],
              "properties": {
                "grantedToV2": {
                  "type": "object",
                  "required": [
                    "group"
                  ],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": [
                        "displayName",
                        "id"
                      ],
                      "properties": {
                        "displayName": {
                          "type": "string",
                          "enum": [
                            "grp1"
                          ]
                        },
                        "id": {
                          "type": "string",
                          "pattern": "^%user_id_pattern%$"
                        }
                      }
                    }
                  }
                },
                "roles": {
                  "type": "array",
                  "minItems": 1,
                  "maxItems": 1,
                  "items": {
                    "type": "string",
                    "pattern": "^%role_id_pattern%$"
                  }
                }
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |


  Scenario Outline: send share invitation for disabled project space to group with different roles
    Given using spaces DAV path
    And user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Admin" has disabled a space "NewSpace"
    When user "Alice" sends the following share invitation for space using the Graph API:
      | space           | NewSpace           |
      | sharee          | grp1               |
      | shareType       | group              |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "404"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "error"
        ],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "message"
            ],
            "properties": {
              "code": {
                "type": "string",
                "enum": ["itemNotFound"]
              },
              "message": {
                "type": "string",
                "pattern": "^stat: error: not found: %user_id_pattern%$"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |


  Scenario Outline: send share invitation for deleted project space to group with different roles
    Given using spaces DAV path
    And user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Admin" has disabled a space "NewSpace"
    And user "Admin" has deleted a space "NewSpace"
    When user "Alice" sends the following share invitation for space using the Graph API:
      | space           | NewSpace           |
      | sharee          | grp1               |
      | shareType       | group              |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "404"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "error"
        ],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "message"
            ],
            "properties": {
              "code": {
                "type": "string",
                "enum": ["itemNotFound"]
              },
              "message": {
                "type": "string",
                "enum": ["stat: error: not found: "]
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |


  Scenario: send share invitation to user for deleted file
    Given user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    And we save it into "FILEID"
    And user "Alice" has deleted file "textfile1.txt"
    When user "Alice" sends the following share invitation with file-id "<<FILEID>>" using the Graph API:
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Viewer   |
    Then the HTTP status code should be "404"
    And for user "Brian" the space Shares should not contain these entries:
      | textfile1.txt |
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "error"
        ],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "message"
            ],
            "properties": {
              "code": {
                "type": "string",
                "enum": ["itemNotFound"]
              },
              "message": {
                "type": "string",
                "enum": ["stat: error: not found: "]
              }
            }
          }
        }
      }
      """


  Scenario: send share invitation to group for deleted file
    Given user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    And we save it into "FILEID"
    And user "Alice" has deleted file "textfile1.txt"
    When user "Alice" sends the following share invitation with file-id "<<FILEID>>" using the Graph API:
      | space           | Personal |
      | sharee          | grp1     |
      | shareType       | group    |
      | permissionsRole | Viewer   |
    Then the HTTP status code should be "404"
    And for user "Brian" the space Shares should not contain these entries:
      | textfile1.txt |
    And for user "Carol" the space Shares should not contain these entries:
      | textfile1.txt |
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "error"
        ],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "message"
            ],
            "properties": {
              "code": {
                "type": "string",
                "enum": ["itemNotFound"]
              },
              "message": {
                "type": "string",
                "enum": ["stat: error: not found: "]
              }
            }
          }
        }
      }
      """


  Scenario Outline: send share invitation for project space resource to user with different roles
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "NewSpace" with content "share space items" to "textfile1.txt"
    And user "Alice" has created a folder "FolderToShare" in space "NewSpace"
    When user "Alice" sends the following share invitation for space using the Graph API:
      | resource        | <resource>         |
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "200"
    And for user "Brian" the space Shares should contain these entries:
      | <resource> |
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "value"
        ],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": [
                "grantedToV2",
                "roles"
              ],
              "properties": {
                "grantedToV2": {
                  "type": "object",
                  "required": [
                    "user"
                  ],
                  "properties": {
                    "user": {
                      "type": "object",
                      "required": [
                        "displayName",
                        "id"
                      ],
                      "properties": {
                        "displayName": {
                          "type": "string",
                          "enum": [
                            "Brian Murphy"
                          ]
                        },
                        "id": {
                          "type": "string",
                          "pattern": "^%user_id_pattern%$"
                        }
                      }
                    }
                  }
                },
                "roles": {
                  "type": "array",
                  "minItems": 1,
                  "maxItems": 1,
                  "items": {
                    "type": "string",
                    "pattern": "^%role_id_pattern%$"
                  }
                }
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource      |
      | Viewer           | textfile1.txt |
      | File Editor      | textfile1.txt |
      | Viewer           | FolderToShare |
      | Editor           | FolderToShare |
      | Uploader         | FolderToShare |


  Scenario Outline: send share invitation for project space resource to group with different roles
    Given using spaces DAV path
    And user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "NewSpace" with content "share space items" to "textfile1.txt"
    And user "Alice" has created a folder "FolderToShare" in space "NewSpace"
    When user "Alice" sends the following share invitation for space using the Graph API:
      | resource        | <resource>         |
      | space           | NewSpace           |
      | sharee          | grp1               |
      | shareType       | group              |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "200"
    And for user "Brian" the space Shares should contain these entries:
      | <resource> |
    And for user "Carol" the space Shares should contain these entries:
      | <resource> |
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "value"
        ],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": [
                "id",
                "roles",
                "grantedToV2"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%permissions_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "minItems": 1,
                  "maxItems": 1,
                  "items": {
                    "type": "string",
                    "pattern": "^%role_id_pattern%$"
                  }
                },
                "grantedToV2": {
                  "type": "object",
                  "required": [
                    "group"
                  ],
                  "properties": {
                    "group": {
                      "type": "object",
                      "required": [
                        "id",
                        "displayName"
                      ],
                      "properties": {
                        "id": {
                          "type": "string",
                          "pattern": "^%group_id_pattern%$"
                        },
                        "displayName": {
                          "type": "string",
                          "enum": [
                            "grp1"
                          ]
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | resource      |
      | Viewer           | textfile1.txt |
      | File Editor      | textfile1.txt |
      | Viewer           | FolderToShare |
      | Editor           | FolderToShare |
      | Uploader         | FolderToShare |

  @issue-8494
  Scenario Outline: try to send share invitation for personal space to user with different roles
    When user "Alice" sends the following share invitation for space using the Graph API:
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "innererror",
              "message"
            ],
            "properties": {
              "code": {
                "const": "invalidRequest"
              },
              "innererror": {
                "type": "object",
                "required": [
                  "date",
                  "request-id"
                ]
              },
              "message": {
                "const": "space type is not eligible for sharing"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |

  @issue-8495
  Scenario Outline: try to share Shares space with a user
    When user "Alice" sends the following share invitation for space using the Graph API:
      | space           | Shares             |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "code",
              "innererror",
              "message"
            ],
            "properties": {
              "code": {
                "const": "invalidRequest"
              },
              "innererror": {
                "type": "object",
                "required": [
                  "date",
                  "request-id"
                ]
              },
              "message": {
                "const": "<error-message>"
              }
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | error-message                                                                                      |
      | Space Viewer     | role not applicable to this resource								      |
      | Space Editor     | role not applicable to this resource                                                               |
      | Manager          | role not applicable to this resource                                                               |
      | Co Owner         | Key: 'DriveItemInvite.Roles' Error:Field validation for 'Roles' failed on the 'available_role' tag |


  Scenario Outline: try to send share invitation with re-sharing permissions
    Given group "grp1" has been created
    And user "Alice" has created folder "FolderToShare"
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
    And user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    When user "Alice" sends the following share invitation using the Graph API:
      | resource          | textfile1.txt        |
      | space             | Personal             |
      | sharee            | grp1                 |
      | shareType         | group                |
      | permissionsAction | <permissions-action> |
    Then the HTTP status code should be "400"
    When user "Alice" sends the following share invitation using the Graph API:
      | resource          | FolderToShare        |
      | space             | Personal             |
      | sharee            | Brian                |
      | shareType         | user                 |
      | permissionsAction | <permissions-action> |
    Then the HTTP status code should be "400"
    Examples:
      | permissions-action |
      | permissions/create |
      | permissions/update |
      | permissions/delete |
      | permissions/deny   |
