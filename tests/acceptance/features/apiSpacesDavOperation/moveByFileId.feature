Feature: moving/renaming file using file id
  As a user
  I want to be able to move or rename files using file id
  So that I can manage my file system

  Background:
    Given using spaces DAV path
    And user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: move a file into a folder inside personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has uploaded file with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" moves a file "/textfile.txt" into "/folder" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "folder" of the space "Personal" should contain these files:
      | textfile.txt |
    But for user "Alice" the space "Personal" should not contain these entries:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file into a sub-folder inside personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has created folder "folder/sub-folder"
    And user "Alice" has uploaded file with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" moves a file "/textfile.txt" into "/folder/sub-folder" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "folder/sub-folder/" of the space "Personal" should contain these files:
      | textfile.txt |
    But for user "Alice" the space "Personal" should not contain these entries:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file from folder to root inside personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has uploaded file with content "some data" to "folder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" moves a file "folder/textfile.txt" into "/" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Personal" should contain these entries:
      | textfile.txt |
    But for user "Alice" folder "folder" of the space "Personal" should not contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file from sub-folder to root inside personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has created folder "folder/sub-folder"
    And user "Alice" has uploaded file with content "some data" to "folder/sub-folder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" moves a file "folder/sub-folder/textfile.txt" into "/" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Personal" should contain these entries:
      | textfile.txt |
    But for user "Alice" folder "folder/sub-folder" of the space "Personal" should not contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |

  @issue-1976
  Scenario Outline: try to move a file into same folder with same name
    And user "Alice" has uploaded file with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    When user "Alice" moves a file "textfile.txt" into "/" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "403"
    And as "Alice" file "textfile.txt" should not exist in the trashbin of the space "Personal"
    And for user "Alice" the content of the file "textfile.txt" of the space "Personal" should be "some data"
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file from personal to share space
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/folder"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folder        |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | <permissions> |
    And user "Brian" has uploaded file with content "some data" to "/test.txt"
    And we save it into "FILEID"
    When user "Brian" moves a file "test.txt" into "folder" inside space "Shares" using file-id path "<dav-path>"
    Then the HTTP status code should be "502"
    And the value of the item "/d:error/s:message" in the response about user "Brian" should be "cross storage moves are not supported, use copy and delete"
    And for user "Brian" folder "/" of the space "Personal" should contain these files:
      | test.txt |
    But for user "Alice" folder "folder" of the space "Personal" should not contain these files:
      | test.txt |
    Examples:
      | permissions | dav-path                          |
      | Editor      | /remote.php/dav/spaces/<<FILEID>> |
      | Editor      | /dav/spaces/<<FILEID>>            |
      | Viewer      | /remote.php/dav/spaces/<<FILEID>> |
      | Viewer      | /dav/spaces/<<FILEID>>            |

  @issue-7618
  Scenario Outline: move a file from personal to project space
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Brian" has uploaded a file inside space "Personal" with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following space share invitation:
      | space           | project-space |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | <space-role>  |
    When user "Brian" moves a file "textfile.txt" into "/" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "<http-status-code>"
    And for user "Brian" folder "/" of the space "Personal" should contain these files:
      | textfile.txt |
    But for user "Brian" folder "/" of the space "project-space" should not contain these files:
      | textfile.txt |
    Examples:
      | space-role   | http-status-code | dav-path                          |
      | Manager      | 502              | /remote.php/dav/spaces/<<FILEID>> |
      | Space Editor | 502              | /remote.php/dav/spaces/<<FILEID>> |
      | Space Viewer | 403              | /remote.php/dav/spaces/<<FILEID>> |
      | Manager      | 502              | /dav/spaces/<<FILEID>>            |
      | Space Editor | 502              | /dav/spaces/<<FILEID>>            |
      | Space Viewer | 403              | /dav/spaces/<<FILEID>>            |

  @issue-7618
  Scenario Outline: move a file to different name from personal space to project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "Personal" with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    When user "Alice" renames a file "/textfile.txt" into "/renamed.txt" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "502"
    And the value of the item "/d:error/s:message" in the response about user "Alice" should be "move:error: not supported: cannot move across spaces"
    And for user "Alice" folder "/" of the space "Personal" should contain these files:
      | textfile.txt |
    But for user "Alice" folder "/" of the space "project-space" should not contain these files:
      | renamed.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file into a folder inside project space (manager/editor)
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has created a folder "/folder" in space "project-space"
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following space share invitation:
      | space           | project-space |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | <space-role>  |
    When user "Brian" moves a file "/textfile.txt" into "/folder" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "folder" of the space "project-space" should contain these files:
      | textfile.txt |
    But for user "Alice" the space "project-space" should not contain these entries:
      | textfile.txt |
    Examples:
      | space-role   | dav-path                          |
      | Manager      | /remote.php/dav/spaces/<<FILEID>> |
      | Space Editor | /remote.php/dav/spaces/<<FILEID>> |
      | Manager      | /dav/spaces/<<FILEID>>            |
      | Space Editor | /dav/spaces/<<FILEID>>            |

  @issue-1976
  Scenario Outline: try to move a file within a project space into a folder with same name
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following space share invitation:
      | space           | project-space |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | <space-role>  |
    When user "Brian" moves a file "textfile.txt" into "/" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "403"
    And as "Alice" file "textfile.txt" should not exist in the trashbin of the space "project-space"
    And for user "Brian" the content of the file "textfile.txt" of the space "project-space" should be "some data"
    Examples:
      | space-role   | dav-path                          |
      | Manager      | /remote.php/dav/spaces/<<FILEID>> |
      | Space Editor | /remote.php/dav/spaces/<<FILEID>> |
      | Space Viewer | /remote.php/dav/spaces/<<FILEID>> |
      | Manager      | /dav/spaces/<<FILEID>>            |
      | Space Viewer | /dav/spaces/<<FILEID>>            |
      | viewer       | /dav/spaces/<<FILEID>>            |


  Scenario Outline: try to move a file into a folder inside project space (viewer)
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has created a folder "/folder" in space "project-space"
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following space share invitation:
      | space           | project-space |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Space Viewer  |
    When user "Brian" moves a file "/textfile.txt" into "/folder" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "403"
    And for user "Alice" the space "project-space" should contain these entries:
      | textfile.txt |
    But for user "Alice" folder "folder" of the space "project-space" should not contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file into a sub-folder inside project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has created a folder "folder/sub-folder" in space "project-space"
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" moves a file "/textfile.txt" into "/folder/sub-folder" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "folder/sub-folder/" of the space "project-space" should contain these files:
      | textfile.txt |
    But for user "Alice" the space "Personal" should not contain these entries:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file from folder to root inside project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has created a folder "folder" in space "project-space"
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "folder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" moves a file "folder/textfile.txt" into "/" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" the space "project-space" should contain these entries:
      | textfile.txt |
    But for user "Alice" folder "folder" of the space "project-space" should not contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file from sub-folder to root inside project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has created a folder "folder/sub-folder" in space "project-space"
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "folder/sub-folder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" moves a file "folder/sub-folder/textfile.txt" into "/" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" the space "project-space" should contain these entries:
      | textfile.txt |
    But for user "Alice" folder "folder/sub-folder" of the space "project-space" should not contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |

  @issue-8116
  Scenario Outline: move a file between two project spaces
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "first-project-space" with the default quota using the Graph API
    And user "Alice" has created a space "second-project-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "first-project-space" with content "first project space" to "file.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following space share invitation:
      | space           | first-project-space |
      | sharee          | Brian               |
      | shareType       | user                |
      | permissionsRole | <from-space-role>   |
    And user "Alice" has sent the following space share invitation:
      | space           | second-project-space |
      | sharee          | Brian                |
      | shareType       | user                 |
      | permissionsRole | <to-space-role>      |
    When user "Brian" moves a file "file.txt" into "/" inside space "second-project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "<http-status-code>"
    And for user "Brian" the space "second-project-space" should not contain these entries:
      | file.txt |
    But for user "Brian" the space "first-project-space" should contain these entries:
      | file.txt |
    Examples:
      | from-space-role | to-space-role | http-status-code | dav-path                          |
      | Manager         | Manager       | 502              | /remote.php/dav/spaces/<<FILEID>> |
      | Space Editor    | Manager       | 502              | /remote.php/dav/spaces/<<FILEID>> |
      | Manager         | Space Editor  | 502              | /remote.php/dav/spaces/<<FILEID>> |
      | Space Editor    | Space Editor  | 502              | /remote.php/dav/spaces/<<FILEID>> |
      | Manager         | Space Viewer  | 403              | /remote.php/dav/spaces/<<FILEID>> |
      | Space Editor    | Space Viewer  | 403              | /remote.php/dav/spaces/<<FILEID>> |
      | Space Viewer    | Manager       | 403              | /remote.php/dav/spaces/<<FILEID>> |
      | Space Viewer    | Space Editor  | 403              | /remote.php/dav/spaces/<<FILEID>> |
      | Space Viewer    | Space Viewer  | 403              | /remote.php/dav/spaces/<<FILEID>> |
      | Manager         | Manager       | 502              | /dav/spaces/<<FILEID>>            |
      | Space Editor    | Manager       | 502              | /dav/spaces/<<FILEID>>            |
      | Manager         | Space Editor  | 502              | /dav/spaces/<<FILEID>>            |
      | Space Editor    | Space Editor  | 502              | /dav/spaces/<<FILEID>>            |
      | Manager         | Space Viewer  | 403              | /dav/spaces/<<FILEID>>            |
      | Space Editor    | Space Viewer  | 403              | /dav/spaces/<<FILEID>>            |
      | Space Viewer    | Manager       | 403              | /dav/spaces/<<FILEID>>            |
      | Space Viewer    | Space Editor  | 403              | /dav/spaces/<<FILEID>>            |
      | Space Viewer    | Space Viewer  | 403              | /dav/spaces/<<FILEID>>            |

  @issue-8116
  Scenario Outline: move a file to different name between project spaces
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "first-project-space" with the default quota using the Graph API
    And user "Alice" has created a space "second-project-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "first-project-space" with content "data from first project space" to "firstProjectSpacetextfile.txt"
    And user "Alice" has uploaded a file inside space "second-project-space" with content "data from second project space" to "secondProjectSpacetextfile.txt"
    And we save it into "FILEID"
    When user "Alice" renames a file "/secondProjectSpacetextfile.txt" into "/renamedSecondProjectSpacetextfile.txt" inside space "first-project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "502"
    And the value of the item "/d:error/s:message" in the response about user "Alice" should be "move:error: not supported: cannot move across spaces"
    And for user "Alice" folder "/" of the space "first-project-space" should contain these files:
      | firstProjectSpacetextfile.txt |
    And for user "Alice" folder "/" of the space "second-project-space" should contain these files:
      | secondProjectSpacetextfile.txt |
    But for user "Alice" the space "first-project-space" should not contain these entries:
      | renamedSecondProjectSpacetextfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file from project to shares space
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following space share invitation:
      | space           | project-space |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | <space-role>  |
    And user "Alice" has created folder "testshare"
    And user "Alice" has sent the following resource share invitation:
      | resource        | testshare     |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | <permissions> |
    When user "Brian" moves a file "textfile.txt" into "testshare" inside space "Shares" using file-id path "<dav-path>"
    Then the HTTP status code should be "502"
    And for user "Brian" folder "/" of the space "project-space" should contain these files:
      | textfile.txt |
    But for user "Brian" folder "testshare" of the space "Shares" should not contain these files:
      | textfile.txt |
    Examples:
      | space-role   | permissions | dav-path                          |
      | Manager      | Editor      | /remote.php/dav/spaces/<<FILEID>> |
      | Space Editor | Editor      | /remote.php/dav/spaces/<<FILEID>> |
      | Space Viewer | Editor      | /remote.php/dav/spaces/<<FILEID>> |
      | Manager      | Viewer      | /remote.php/dav/spaces/<<FILEID>> |
      | Space Editor | Viewer      | /remote.php/dav/spaces/<<FILEID>> |
      | Space Viewer | Viewer      | /remote.php/dav/spaces/<<FILEID>> |
      | Manager      | Editor      | /dav/spaces/<<FILEID>>            |
      | Space Editor | Editor      | /dav/spaces/<<FILEID>>            |
      | Space Viewer | Editor      | /dav/spaces/<<FILEID>>            |
      | Manager      | Viewer      | /dav/spaces/<<FILEID>>            |
      | Space Editor | Viewer      | /dav/spaces/<<FILEID>>            |
      | Space Viewer | Viewer      | /dav/spaces/<<FILEID>>            |

  @issue-7618
  Scenario Outline: move a file from project to personal space
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following space share invitation:
      | space           | project-space |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | <space-role>  |
    When user "Brian" moves a file "/textfile.txt" into "/" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "<http-status-code>"
    And for user "Brian" folder "/" of the space "project-space" should contain these files:
      | textfile.txt |
    But for user "Brian" folder "/" of the space "Personal" should not contain these files:
      | textfile.txt |
    Examples:
      | space-role   | http-status-code | dav-path                          |
      | Manager      | 502              | /remote.php/dav/spaces/<<FILEID>> |
      | Space Editor | 502              | /remote.php/dav/spaces/<<FILEID>> |
      | Space Viewer | 403              | /remote.php/dav/spaces/<<FILEID>> |
      | Manager      | 502              | /dav/spaces/<<FILEID>>            |
      | Space Editor | 502              | /dav/spaces/<<FILEID>>            |
      | Space Viewer | 403              | /dav/spaces/<<FILEID>>            |

  @issue-7618
  Scenario Outline: move a file to different name from project space to personal space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    When user "Alice" renames a file "/textfile.txt" into "/renamed.txt" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "502"
    And the value of the item "/d:error/s:message" in the response about user "Alice" should be "move:error: not supported: cannot move across spaces"
    And for user "Alice" folder "/" of the space "project-space" should contain these files:
      | textfile.txt |
    But for user "Alice" folder "/" of the space "Personal" should not contain these files:
      | renamed.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |

  @issue-7617
  Scenario Outline: move a file into a folder within a shared folder (edit permissions)
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "folder"
    And user "Alice" has created folder "folder/sub-folder"
    And user "Alice" has uploaded file with content "some data" to "folder/test.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folder        |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | <permissions> |
    When user "Brian" moves a file "Shares/folder/test.txt" into "folder/sub-folder" inside space "Shares" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Brian" folder "folder/sub-folder" of the space "Shares" should contain these files:
      | test.txt |
    And for user "Alice" folder "folder/sub-folder" of the space "Personal" should contain these files:
      | test.txt |
    But for user "Brian" folder "folder" of the space "Shares" should not contain these files:
      | test.txt |
    And for user "Alice" folder "folder" of the space "Personal" should not contain these files:
      | test.txt |
    Examples:
      | permissions | dav-path                          |
      | Editor      | /remote.php/dav/spaces/<<FILEID>> |
      | Editor      | /dav/spaces/<<FILEID>>            |

  @issue-1976
  Scenario Outline: sharee tries to move a file into same shared folder with same name
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "folder"
    And user "Alice" has uploaded file with content "some data" to "folder/test.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folder        |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | <permissions> |
    When user "Brian" moves a file "Shares/folder/test.txt" into "folder" inside space "Shares" using file-id path "<dav-path>"
    Then the HTTP status code should be "403"
    And as "Alice" file "test.txt" should not exist in the trashbin of the space "Personal"
    And for user "Brian" the content of the file "folder/test.txt" of the space "Shares" should be "some data"
    And for user "Alice" the content of the file "folder/test.txt" of the space "Personal" should be "some data"
    Examples:
      | permissions | dav-path                          |
      | Editor      | /remote.php/dav/spaces/<<FILEID>> |
      | Editor      | /dav/spaces/<<FILEID>>            |
      | Viewer      | /remote.php/dav/spaces/<<FILEID>> |
      | Viewer      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: try to move a file into a folder within a shared folder (read permissions)
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "folder"
    And user "Alice" has created folder "folder/sub-folder"
    And user "Alice" has uploaded file with content "some data" to "folder/test.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folder   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Viewer   |
    When user "Brian" moves a file "Shares/folder/test.txt" into "folder/sub-folder" inside space "Shares" using file-id path "<dav-path>"
    Then the HTTP status code should be "502"
    And for user "Brian" folder "folder/sub-folder" of the space "Shares" should not contain these files:
      | test.txt |
    And for user "Alice" folder "folder/sub-folder" of the space "Personal" should not contain these files:
      | test.txt |
    But for user "Brian" folder "folder" of the space "Shares" should contain these files:
      | test.txt |
    And for user "Alice" folder "folder" of the space "Personal" should contain these files:
      | test.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file from one shared folder to another shared folder
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "testshare1"
    And user "Alice" has created folder "testshare2"
    And user "Alice" has uploaded file with content "some data" to "testshare1/textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following resource share invitation:
      | resource        | testshare1         |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <from-permissions> |
    And user "Alice" has sent the following resource share invitation:
      | resource        | testshare2       |
      | space           | Personal         |
      | sharee          | Brian            |
      | shareType       | user             |
      | permissionsRole | <to-permissions> |
    When user "Brian" moves a file "Shares/testshare1/textfile.txt" into "testshare2" inside space "Shares" using file-id path "<dav-path>"
    Then the HTTP status code should be "502"
    And for user "Brian" folder "testshare1" of the space "Shares" should contain these files:
      | textfile.txt |
    But for user "Brian" folder "testshare2" of the space "Shares" should not contain these files:
      | textfile.txt |
    Examples:
      | from-permissions | to-permissions | dav-path                          |
      | Editor           | Editor         | /remote.php/dav/spaces/<<FILEID>> |
      | Editor           | Viewer         | /remote.php/dav/spaces/<<FILEID>> |
      | Viewer           | Editor         | /remote.php/dav/spaces/<<FILEID>> |
      | Viewer           | Viewer         | /remote.php/dav/spaces/<<FILEID>> |
      | Editor           | Editor         | /dav/spaces/<<FILEID>>            |
      | Editor           | Viewer         | /dav/spaces/<<FILEID>>            |
      | Viewer           | Editor         | /dav/spaces/<<FILEID>>            |
      | Viewer           | Viewer         | /dav/spaces/<<FILEID>>            |

  @issue-8124
  Scenario Outline: move a file from share to personal space
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/folder"
    And user "Alice" has uploaded file with content "some data" to "/folder/test.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folder        |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | <permissions> |
    When user "Brian" moves a file "Shares/folder/test.txt" into "/" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "<http-status-code>"
    And for user "Brian" folder "folder" of the space "Shares" should contain these files:
      | test.txt |
    And for user "Brian" folder "/" of the space "Personal" should not contain these files:
      | test.txt |
    Examples:
      | permissions | dav-path                          | http-status-code |
      | Editor      | /remote.php/dav/spaces/<<FILEID>> | 502              |
      | Editor      | /dav/spaces/<<FILEID>>            | 502              |
      | Viewer      | /remote.php/dav/spaces/<<FILEID>> | 403              |
      | Viewer      | /dav/spaces/<<FILEID>>            | 403              |

  @issue-8125
  Scenario Outline: move a file from shares to project space
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has sent the following space share invitation:
      | space           | project-space |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | <space-role>  |
    And user "Alice" has created folder "testshare"
    And user "Alice" has uploaded file with content "some data" to "testshare/textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following resource share invitation:
      | resource        | testshare     |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | <permissions> |
    When user "Brian" moves a file "Shares/testshare/textfile.txt" into "/" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "<http-status-code>"
    And for user "Brian" folder "testshare" of the space "Shares" should contain these files:
      | textfile.txt |
    But for user "Brian" folder "/" of the space "project-space" should not contain these files:
      | textfile.txt |
    Examples:
      | space-role   | permissions | http-status-code | dav-path                          |
      | Manager      | Editor      | 502              | /remote.php/dav/spaces/<<FILEID>> |
      | Space Editor | Editor      | 502              | /remote.php/dav/spaces/<<FILEID>> |
      | Space Viewer | Editor      | 403              | /remote.php/dav/spaces/<<FILEID>> |
      | Manager      | Viewer      | 403              | /remote.php/dav/spaces/<<FILEID>> |
      | Space Editor | Viewer      | 403              | /remote.php/dav/spaces/<<FILEID>> |
      | Space Viewer | Viewer      | 403              | /remote.php/dav/spaces/<<FILEID>> |
      | Manager      | Editor      | 502              | /dav/spaces/<<FILEID>>            |
      | Space Editor | Editor      | 502              | /dav/spaces/<<FILEID>>            |
      | Space Viewer | Editor      | 403              | /dav/spaces/<<FILEID>>            |
      | Manager      | Viewer      | 403              | /dav/spaces/<<FILEID>>            |
      | Space Editor | Viewer      | 403              | /dav/spaces/<<FILEID>>            |
      | Space Viewer | Viewer      | 403              | /dav/spaces/<<FILEID>>            |


  Scenario Outline: rename a root file inside personal space
    Given user "Alice" has uploaded file with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    When user "Alice" renames a file "textfile.txt" into "renamed.txt" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Personal" should contain these entries:
      | renamed.txt |
    But for user "Alice" the space "Personal" should not contain these entries:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: rename a file and move into a folder inside personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has uploaded file with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" renames a file "textfile.txt" into "/folder/renamed.txt" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "folder" of the space "Personal" should contain these files:
      | renamed.txt |
    But for user "Alice" the space "Personal" should not contain these entries:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: rename a file and move into a sub-folder inside personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has created folder "folder/sub-folder"
    And user "Alice" has uploaded file with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" renames a file "textfile.txt" into "/folder/sub-folder/renamed.txt" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "folder/sub-folder" of the space "Personal" should contain these files:
      | renamed.txt |
    But for user "Alice" the space "Personal" should not contain these entries:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: rename a file and move from a folder to root inside personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has uploaded file with content "some data" to "folder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" renames a file "folder/textfile.txt" into "/renamed.txt" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Personal" should contain these entries:
      | renamed.txt |
    But for user "Alice" folder "folder" of the space "Personal" should not contain these files:
      | renamed.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: rename a file and move from sub-folder to root inside personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has created folder "folder/sub-folder"
    And user "Alice" has uploaded file with content "some data" to "folder/sub-folder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" renames a file "folder/sub-folder/textfile.txt" into "/renamed.txt" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Personal" should contain these files:
      | renamed.txt |
    But for user "Alice" folder "folder/sub-folder" of the space "Personal" should not contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |

  @issue-7617
  Scenario Outline: move a file to a different name into a sub-folder inside share space (editor permissions)
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/folder"
    And user "Alice" has created folder "/folder/sub-folder"
    And user "Alice" has uploaded file with content "some data" to "/folder/test.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folder   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    When user "Brian" renames a file "Shares/folder/test.txt" into "folder/sub-folder/renamed.txt" inside space "Shares" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Brian" folder "folder/sub-folder" of the space "Shares" should contain these files:
      | renamed.txt |
    But for user "Brian" folder "folder" of the space "Shares" should not contain these files:
      | test.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file to a different name into a sub-folder inside share space (read permissions)
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/folder"
    And user "Alice" has created folder "/folder/sub-folder"
    And user "Alice" has uploaded file with content "some data" to "/folder/test.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folder   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Viewer   |
    When user "Brian" renames a file "Shares/folder/test.txt" into "folder/sub-folder/renamed.txt" inside space "Shares" using file-id path "<dav-path>"
    Then the HTTP status code should be "502"
    And for user "Brian" folder "folder" of the space "Shares" should contain these files:
      | test.txt |
    But for user "Brian" folder "folder/sub-folder" of the space "Shares" should not contain these files:
      | renamed.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |

  @issue-6739
  Scenario Outline: try to move personal file to other spaces using its id as the destination
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "myspace" with the default quota using the Graph API
    And user "Alice" has uploaded file with content "some data" to "textfile.txt"
    When user "Alice" tries to move file "textfile.txt" of space "Personal" to space "<space-name>" using its id in destination path "<dav-path>"
    Then the HTTP status code should be "<http-status-code>"
    And for user "Alice" the space "Personal" should contain these entries:
      | textfile.txt |
    Examples:
      | dav-path               | space-name | http-status-code |
      | /remote.php/dav/spaces | Personal   | 409              |
      | /dav/spaces            | Personal   | 409              |
      | /remote.php/dav/spaces | myspace    | 400              |
      | /dav/spaces            | myspace    | 400              |
      | /remote.php/dav/spaces | Shares     | 404              |
      | /dav/spaces            | Shares     | 404              |

  @issue-6739
  Scenario Outline: try to move project file to other spaces using its id as the destination
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "myspace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "myspace" with content "some data" to "textfile.txt"
    When user "Alice" tries to move file "textfile.txt" of space "myspace" to space "<space-name>" using its id in destination path "<dav-path>"
    Then the HTTP status code should be "<http-status-code>"
    And for user "Alice" the space "myspace" should contain these entries:
      | textfile.txt |
    Examples:
      | dav-path               | space-name | http-status-code |
      | /remote.php/dav/spaces | Personal   | 400              |
      | /dav/spaces            | Personal   | 400              |
      | /remote.php/dav/spaces | myspace    | 409              |
      | /dav/spaces            | myspace    | 409              |
      | /remote.php/dav/spaces | Shares     | 404              |
      | /dav/spaces            | Shares     | 404              |

  @issue-6739
  Scenario Outline: move a file to folder using its id as the destination (Personal space)
    Given user "Alice" has uploaded file with content "some data" to "textfile.txt"
    And user "Alice" has created folder "docs"
    When user "Alice" moves file "textfile.txt" of space "Personal" to folder "docs" using its id in destination path "<dav-path>"
    Then the HTTP status code should be "204"
    And the content of file "docs" for user "Alice" should be "some data"
    And as "Alice" file "textfile.txt" should not exist
    And as "Alice" folder "docs" should not exist
    And as "Alice" folder "docs" should exist in the trashbin of the space "Personal"
    Examples:
      | dav-path               |
      | /remote.php/dav/spaces |
      | /dav/spaces            |

  @issue-6739
  Scenario Outline: move a file to folder using its id as the destination (Project space)
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "myspace" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "myspace" with content "some data" to "textfile.txt"
    And user "Alice" has created a folder "docs" in space "myspace"
    When user "Alice" moves file "textfile.txt" of space "myspace" to folder "docs" using its id in destination path "<dav-path>"
    Then the HTTP status code should be "204"
    And for user "Alice" the content of the file "docs" of the space "myspace" should be "some data"
    And as "Alice" folder "docs" should exist in the trashbin of the space "myspace"
    And for user "Alice" folder "/" of the space "myspace" should not contain these files:
      | textfile.txt |
    Examples:
      | dav-path               |
      | /remote.php/dav/spaces |
      | /dav/spaces            |
