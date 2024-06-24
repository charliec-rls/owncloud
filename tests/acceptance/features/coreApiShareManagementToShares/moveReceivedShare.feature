@skipOnReva @issue-1289 @issue-1328
Feature: sharing

  Background:
    Given using OCS API version "1"
    And these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Carol    |

  @issue-8242
  Scenario Outline: sharer renames the shared item (old/new webdav)
    Given user "Alice" has uploaded file with content "foo" to "sharefile.txt"
    And using <dav-path-version> DAV path
    And user "Alice" has sent the following resource share invitation:
      | resource        | sharefile.txt |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    And user "Alice" has sent the following resource share invitation:
      | resource        | sharefile.txt |
      | space           | Personal      |
      | sharee          | Carol         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    When user "Alice" moves file "sharefile.txt" to "renamedsharefile.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Alice" file "renamedsharefile.txt" should exist
    And as "Brian" file "Shares/sharefile.txt" should exist
    And as "Carol" file "Shares/sharefile.txt" should exist
    When user "Alice" sends HTTP method "PROPFIND" to URL "<dav-path>"
    Then the HTTP status code should be "207"
    And as user "Alice" the value of the item "//oc:name" of path "<dav-path>/renamedsharefile.txt" in the response should be "renamedsharefile.txt"
    And as user "Alice" the value of the item "//d:displayname" of path "<dav-path>/renamedsharefile.txt" in the response should be "renamedsharefile.txt"
    When user "Brian" sends HTTP method "PROPFIND" to URL "<dav-path>/Shares"
    Then the HTTP status code should be "207"
    And as user "Brian" the value of the item "//oc:name" of path "<dav-path>/Shares/sharefile.txt" in the response should be "sharefile.txt"
    And as user "Brian" the value of the item "//d:displayname" of path "<dav-path>/Shares/sharefile.txt" in the response should be "sharefile.txt"
    When user "Carol" sends HTTP method "PROPFIND" to URL "<dav-path>/Shares"
    Then the HTTP status code should be "207"
    And as user "Carol" the value of the item "//oc:name" of path "<dav-path>/Shares/sharefile.txt" in the response should be "sharefile.txt"
    And as user "Carol" the value of the item "//d:displayname" of path "<dav-path>/Shares/sharefile.txt" in the response should be "sharefile.txt"
    Examples:
      | dav-path-version | dav-path                         |
      | old              | /remote.php/webdav               |
      | new              | /remote.php/dav/files/%username% |
      | new              | /dav/files/%username%            |

  @issue-8242
  Scenario Outline: sharer renames the shared item (spaces webdav)
    Given user "Alice" has uploaded file with content "foo" to "sharefile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | sharefile.txt |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    And user "Alice" has sent the following resource share invitation:
      | resource        | sharefile.txt |
      | space           | Personal      |
      | sharee          | Carol         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    When user "Alice" moves file "sharefile.txt" to "renamedsharefile.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Alice" file "renamedsharefile.txt" should exist
    And as "Brian" file "Shares/sharefile.txt" should exist
    And as "Carol" file "Shares/sharefile.txt" should exist
    And using spaces DAV path
    When user "Alice" sends HTTP method "PROPFIND" to URL "<dav-path-personal>"
    Then the HTTP status code should be "207"
    And as user "Alice" the value of the item "//oc:name" of path "<dav-path-personal>/renamedsharefile.txt" in the response should be "renamedsharefile.txt"
    And as user "Alice" the value of the item "//d:displayname" of path "<dav-path-personal>/renamedsharefile.txt" in the response should be "renamedsharefile.txt"
    When user "Brian" sends HTTP method "PROPFIND" to URL "<dav-path>"
    Then the HTTP status code should be "207"
    And as user "Brian" the value of the item "//oc:name" of path "<dav-path>/sharefile.txt" in the response should be "sharefile.txt"
    And as user "Brian" the value of the item "//d:displayname" of path "<dav-path>/sharefile.txt" in the response should be "sharefile.txt"
    When user "Carol" sends HTTP method "PROPFIND" to URL "<dav-path>"
    Then the HTTP status code should be "207"
    And as user "Carol" the value of the item "//oc:name" of path "<dav-path>/sharefile.txt" in the response should be "sharefile.txt"
    And as user "Carol" the value of the item "//d:displayname" of path "<dav-path>/sharefile.txt" in the response should be "sharefile.txt"
    Examples:
      | dav-path                                 | dav-path-personal                |
      | /remote.php/dav/spaces/%shares_drive_id% | /remote.php/dav/spaces/%spaceid% |
      | /dav/spaces/%shares_drive_id%            | /remote.php/dav/spaces/%spaceid% |

  @issue-8242
  Scenario Outline: share receiver renames the shared item (old/new webdav)
    Given user "Alice" has uploaded file with content "foo" to "/sharefile.txt"
    And using <dav-path-version> DAV path
    And user "Alice" has sent the following resource share invitation:
      | resource        | sharefile.txt      |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    And user "Alice" has sent the following resource share invitation:
      | resource        | sharefile.txt      |
      | space           | Personal           |
      | sharee          | Carol              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Carol" moves file "Shares/sharefile.txt" to "Shares/renamedsharefile.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Carol" file "Shares/renamedsharefile.txt" should exist
    And as "Brian" file "Shares/sharefile.txt" should exist
    And as "Alice" file "sharefile.txt" should exist
    When user "Carol" sends HTTP method "PROPFIND" to URL "<dav-path>/Shares"
    Then the HTTP status code should be "207"
    And as user "Carol" the value of the item "//oc:name" of path "<dav-path>/Shares/renamedsharefile.txt" in the response should be "renamedsharefile.txt"
    And as user "Carol" the value of the item "//d:displayname" of path "<dav-path>/Shares/renamedsharefile.txt" in the response should be "renamedsharefile.txt"
    When user "Alice" sends HTTP method "PROPFIND" to URL "<dav-path>"
    Then the HTTP status code should be "207"
    And as user "Alice" the value of the item "//oc:name" of path "<dav-path>/sharefile.txt" in the response should be "sharefile.txt"
    And as user "Alice" the value of the item "//d:displayname" of path "<dav-path>/sharefile.txt" in the response should be "sharefile.txt"
    When user "Brian" sends HTTP method "PROPFIND" to URL "<dav-path>/Shares"
    Then the HTTP status code should be "207"
    And as user "Brian" the value of the item "//oc:name" of path "<dav-path>/Shares/sharefile.txt" in the response should be "sharefile.txt"
    And as user "Brian" the value of the item "//d:displayname" of path "<dav-path>/Shares/sharefile.txt" in the response should be "sharefile.txt"
    Examples:
      | dav-path-version | dav-path                         | permissions-role |
      | old              | /remote.php/webdav               | Viewer           |
      | new              | /remote.php/dav/files/%username% | Viewer           |
      | new              | /dav/files/%username%            | Viewer           |
      | old              | /remote.php/webdav               | Secure viewer    |
      | new              | /remote.php/dav/files/%username% | Secure viewer    |
      | new              | /dav/files/%username%            | Secure viewer    |

  @issue-8242
  Scenario Outline: share receiver renames the shared item (spaces webdav)
    Given user "Alice" has uploaded file with content "foo" to "/sharefile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | sharefile.txt      |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    And user "Alice" has sent the following resource share invitation:
      | resource        | sharefile.txt      |
      | space           | Personal           |
      | sharee          | Carol              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Carol" moves file "Shares/sharefile.txt" to "Shares/renamedsharefile.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Carol" file "Shares/renamedsharefile.txt" should exist
    And as "Brian" file "Shares/sharefile.txt" should exist
    And as "Alice" file "sharefile.txt" should exist
    And using spaces DAV path
    When user "Carol" sends HTTP method "PROPFIND" to URL "<dav-path>"
    Then the HTTP status code should be "207"
    And as user "Carol" the value of the item "//oc:name" of path "<dav-path>/renamedsharefile.txt" in the response should be "renamedsharefile.txt"
    And as user "Carol" the value of the item "//d:displayname" of path "<dav-path>/renamedsharefile.txt" in the response should be "renamedsharefile.txt"
    When user "Alice" sends HTTP method "PROPFIND" to URL "<dav-path-personal>"
    Then the HTTP status code should be "207"
    And as user "Alice" the value of the item "//oc:name" of path "<dav-path-personal>/sharefile.txt" in the response should be "sharefile.txt"
    And as user "Alice" the value of the item "//d:displayname" of path "<dav-path-personal>/sharefile.txt" in the response should be "sharefile.txt"
    When user "Brian" sends HTTP method "PROPFIND" to URL "<dav-path>"
    Then the HTTP status code should be "207"
    And as user "Brian" the value of the item "//oc:name" of path "<dav-path>/sharefile.txt" in the response should be "sharefile.txt"
    And as user "Brian" the value of the item "//d:displayname" of path "<dav-path>/sharefile.txt" in the response should be "sharefile.txt"
    Examples:
      | dav-path                                 | dav-path-personal                | permissions-role |
      | /remote.php/dav/spaces/%shares_drive_id% | /remote.php/dav/spaces/%spaceid% | Viewer           |
      | /dav/spaces/%shares_drive_id%            | /remote.php/dav/spaces/%spaceid% | Viewer           |
      | /remote.php/dav/spaces/%shares_drive_id% | /remote.php/dav/spaces/%spaceid% | Secure viewer    |
      | /dav/spaces/%shares_drive_id%            | /remote.php/dav/spaces/%spaceid% | Secure viewer    |


  Scenario: keep group share when the one user renames the share and the user is deleted
    Given group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Carol" has been added to group "grp1"
    And user "Alice" has created folder "/TMP"
    And user "Alice" has sent the following resource share invitation:
      | resource        | TMP      |
      | space           | Personal |
      | sharee          | grp1     |
      | shareType       | group    |
      | permissionsRole | Viewer   |
    When user "Carol" moves folder "/Shares/TMP" to "/Shares/new" using the WebDAV API
    Then the HTTP status code should be "201"
    When the administrator deletes user "Carol" using the provisioning API
    Then the HTTP status code should be "204"
    And as "Alice" folder "Shares/TMP" should not exist


  Scenario: receiver renames a received share with read, change permissions inside the Shares folder
    Given user "Alice" has created folder "folderToShare"
    And user "Alice" has uploaded file with content "thisIsAFileInsideTheSharedFolder" to "/folderToShare/fileInside"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folderToShare |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Editor        |
    When user "Brian" moves folder "/Shares/folderToShare" to "/Shares/myFolder" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Brian" folder "/Shares/myFolder" should exist
    But as "Alice" folder "/Shares/myFolder" should not exist
    When user "Brian" moves file "/Shares/myFolder/fileInside" to "/Shares/myFolder/renamedFile" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Brian" file "/Shares/myFolder/renamedFile" should exist
    And as "Alice" file "/folderToShare/renamedFile" should exist
    But as "Alice" file "/folderToShare/fileInside" should not exist


  Scenario Outline: receiver tries to rename a received share with read permissions inside the Shares folder
    Given user "Alice" has created folder "folderToShare"
    And user "Alice" has created folder "folderToShare/folderInside"
    And user "Alice" has uploaded file with content "thisIsAFileInsideTheSharedFolder" to "/folderToShare/fileInside"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folderToShare      |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Brian" moves folder "/Shares/folderToShare" to "/Shares/myFolder" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Brian" folder "/Shares/myFolder" should exist
    But as "Alice" folder "/Shares/myFolder" should not exist
    When user "Brian" moves file "/Shares/myFolder/fileInside" to "/Shares/myFolder/renamedFile" using the WebDAV API
    Then the HTTP status code should be "403"
    And as "Brian" file "/Shares/myFolder/renamedFile" should not exist
    But as "Brian" file "Shares/myFolder/fileInside" should exist
    When user "Brian" moves folder "/Shares/myFolder/folderInside" to "/Shares/myFolder/renamedFolder" using the WebDAV API
    Then the HTTP status code should be "403"
    And as "Brian" folder "/Shares/myFolder/renamedFolder" should not exist
    But as "Brian" folder "Shares/myFolder/folderInside" should exist
    Examples:
      | permissions-role |
      | Viewer           |
      | Secure viewer    |


  Scenario: receiver renames a received folder share to a different name on the same folder
    Given user "Alice" has created folder "PARENT"
    And user "Alice" has sent the following resource share invitation:
      | resource        | PARENT   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    When user "Brian" moves folder "/Shares/PARENT" to "/Shares/myFolder" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Brian" folder "/Shares/myFolder" should exist
    But as "Alice" folder "myFolder" should not exist


  Scenario: receiver renames a received file share to different name on the same folder
    Given user "Alice" has uploaded file "filesForUpload/textfile.txt" to "fileToShare.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | fileToShare.txt |
      | space           | Personal        |
      | sharee          | Brian           |
      | shareType       | user            |
      | permissionsRole | File Editor     |
    When user "Brian" moves file "/Shares/fileToShare.txt" to "/Shares/newFile.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Brian" file "/Shares/newFile.txt" should exist
    But as "Alice" file "newFile.txt" should not exist


  Scenario: receiver renames a received file share to different name on the same folder for group sharing
    Given group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "fileToShare.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | fileToShare.txt |
      | space           | Personal        |
      | sharee          | grp1            |
      | shareType       | group           |
      | permissionsRole | File Editor     |
    When user "Brian" moves file "/Shares/fileToShare.txt" to "/Shares/newFile.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Brian" file "/Shares/newFile.txt" should exist
    But as "Alice" file "newFile.txt" should not exist


  Scenario: receiver renames a received folder share to different name on the same folder for group sharing
    Given group "grp1" has been created
    And user "Alice" has created folder "PARENT"
    And user "Brian" has been added to group "grp1"
    And user "Alice" has sent the following resource share invitation:
      | resource        | PARENT   |
      | space           | Personal |
      | sharee          | grp1     |
      | shareType       | group    |
      | permissionsRole | Viewer   |
    When user "Brian" moves folder "/Shares/PARENT" to "/Shares/myFolder" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Brian" folder "/Shares/myFolder" should exist
    But as "Alice" folder "myFolder" should not exist


  Scenario: receiver renames a received file share with read,update permissions inside the Shares folder in group sharing
    Given group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "fileToShare.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | fileToShare.txt |
      | space           | Personal        |
      | sharee          | grp1            |
      | shareType       | group           |
      | permissionsRole | File Editor     |
    When user "Brian" moves folder "/Shares/fileToShare.txt" to "/Shares/newFile.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Brian" file "/Shares/newFile.txt" should exist
    But as "Alice" file "/Shares/newFile.txt" should not exist


  Scenario: receiver renames a received folder share with read, change permissions inside the Shares folder in group sharing
    Given group "grp1" has been created
    And user "Alice" has created folder "PARENT"
    And user "Brian" has been added to group "grp1"
    And user "Alice" has sent the following resource share invitation:
      | resource        | PARENT   |
      | space           | Personal |
      | sharee          | grp1     |
      | shareType       | group    |
      | permissionsRole | Editor   |
    When user "Brian" moves folder "/Shares/PARENT" to "/Shares/myFolder" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Brian" folder "/Shares/myFolder" should exist
    But as "Alice" folder "/Shares/myFolder" should not exist


  Scenario: receiver renames a received file share with share, read permissions inside the Shares folder in group sharing)
    Given group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "fileToShare.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | fileToShare.txt |
      | space           | Personal        |
      | sharee          | grp1            |
      | shareType       | group           |
      | permissionsRole | Viewer          |
    When user "Brian" moves file "/Shares/fileToShare.txt" to "/Shares/newFile.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Brian" file "/Shares/newFile.txt" should exist
    But as "Alice" file "/Shares/newFile.txt" should not exist


  Scenario: receiver renames a received folder share with read permissions inside the Shares folder in group sharing
    Given group "grp1" has been created
    And user "Alice" has created folder "PARENT"
    And user "Brian" has been added to group "grp1"
    And user "Alice" has sent the following resource share invitation:
      | resource        | PARENT   |
      | space           | Personal |
      | sharee          | grp1     |
      | shareType       | group    |
      | permissionsRole | Viewer   |
    When user "Brian" moves folder "/Shares/PARENT" to "/Shares/myFolder" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Brian" folder "/Shares/myFolder" should exist
    But as "Alice" folder "/Shares/myFolder" should not exist

  @issue-2141
  Scenario Outline: receiver renames a received folder share to name with special characters in group sharing
    Given group "grp1" has been created
    And user "Carol" has been added to group "grp1"
    And user "Alice" has created folder "<sharer-folder>"
    And user "Alice" has created folder "<group-folder>"
    And user "Alice" has sent the following resource share invitation:
      | resource        | <sharer-folder> |
      | space           | Personal        |
      | sharee          | Brian           |
      | shareType       | user            |
      | permissionsRole | Editor          |
    When user "Brian" moves folder "/Shares/<sharer-folder>" to "/Shares/<receiver-folder>" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Alice" folder "<receiver-folder>" should not exist
    And as "Brian" folder "/Shares/<receiver-folder>" should exist
    When user "Alice" shares folder "<group-folder>" with group "grp1" using the sharing API
    And user "Carol" moves folder "/Shares/<group-folder>" to "/Shares/<receiver-folder>" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Alice" folder "<receiver-folder>" should not exist
    But as "Carol" folder "/Shares/<receiver-folder>" should exist
    Examples:
      | sharer-folder | group-folder    | receiver-folder |
      | ?abc=oc #     | ?abc=oc g%rp#   | # oc?test=oc&a  |
      | @a#8a=b?c=d   | @a#8a=b?c=d grp | ?a#8 a=b?c=d    |

  @issue-2141
  Scenario Outline: receiver renames a received file share to name with special characters with read, change permissions in group sharing
    Given group "grp1" has been created
    And user "Carol" has been added to group "grp1"
    And user "Alice" has created folder "<sharer-folder>"
    And user "Alice" has created folder "<group-folder>"
    And user "Alice" has uploaded file with content "thisIsAFileInsideTheSharedFolder" to "/<sharer-folder>/fileInside"
    And user "Alice" has uploaded file with content "thisIsAFileInsideTheSharedFolder" to "/<group-folder>/fileInside"
    And user "Alice" has sent the following resource share invitation:
      | resource        | <sharer-folder> |
      | space           | Personal        |
      | sharee          | Brian           |
      | shareType       | user            |
      | permissionsRole | Editor          |
    When user "Brian" moves folder "/Shares/<sharer-folder>/fileInside" to "/Shares/<sharer-folder>/<receiver_file>" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Alice" file "<sharer-folder>/<receiver_file>" should exist
    And as "Brian" file "/Shares/<sharer-folder>/<receiver_file>" should exist
    When user "Alice" shares folder "<group-folder>" with group "grp1" with permissions "read,change" using the sharing API
    And user "Carol" moves folder "/Shares/<group-folder>/fileInside" to "/Shares/<group-folder>/<receiver_file>" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Alice" file "<group-folder>/<receiver_file>" should exist
    And as "Carol" file "/Shares/<group-folder>/<receiver_file>" should exist
    Examples:
      | sharer-folder | group-folder    | receiver_file  |
      | ?abc=oc #     | ?abc=oc g%rp#   | # oc?test=oc&a |
      | @a#8a=b?c=d   | @a#8a=b?c=d grp | ?a#8 a=b?c=d   |
