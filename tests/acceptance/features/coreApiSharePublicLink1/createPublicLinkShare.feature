@api @issue-1269 @issue-1293

Feature: create a public link share
  As a user
  I want to create public links
  So that I can share resources to people who aren't owncloud users

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files

  @smokeTest
  Scenario Outline: new public link share of a file using the default permissions only grants read access using the public WebDAV API
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has uploaded file with content "Random data" to "/randomfile.txt"
    When user "Alice" creates a public link share using the sharing API with settings
      | path | randomfile.txt |
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And the fields of the last response to user "Alice" should include
      | item_type              | file            |
      | mimetype               | text/plain      |
      | file_target            | /randomfile.txt |
      | path                   | /randomfile.txt |
      | permissions            | read            |
      | share_type             | public_link     |
      | displayname_file_owner | %displayname%   |
      | displayname_owner      | %displayname%   |
      | uid_file_owner         | %username%      |
      | uid_owner              | %username%      |
      | name                   |                 |
    And the public should be able to download the last publicly shared file using the old public WebDAV API without a password and the content should be "Random data"
    And the public should be able to download the last publicly shared file using the new public WebDAV API without a password and the content should be "Random data"
    And the public upload to the last publicly shared file using the old public WebDAV API should fail with HTTP status code "403"
    And the public upload to the last publicly shared file using the new public WebDAV API should fail with HTTP status code "403"
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |

  @smokeTest
  Scenario Outline: creating a new public link share of a file with password using the new public WebDAV API
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has uploaded file with content "Random data" to "/randomfile.txt"
    When user "Alice" creates a public link share using the sharing API with settings
      | path     | randomfile.txt |
      | password | %public%       |
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And the fields of the last response to user "Alice" should include
      | item_type              | file            |
      | mimetype               | text/plain      |
      | file_target            | /randomfile.txt |
      | path                   | /randomfile.txt |
      | permissions            | read            |
      | share_type             | public_link     |
      | displayname_file_owner | %displayname%   |
      | displayname_owner      | %displayname%   |
      | uid_file_owner         | %username%      |
      | uid_owner              | %username%      |
      | name                   |                 |
    And the public should be able to download the last publicly shared file using the new public WebDAV API with password "%public%" and the content should be "Random data"
    And the HTTP status code should be "200"
    And the public download of the last publicly shared file using the new public WebDAV API with password "%regular%" should fail with HTTP status code "401"
    And the value of the item "//s:message" in the response should match "/Username or password was incorrect/"
    And the public download of the last publicly shared file using the new public WebDAV API without a password should fail with HTTP status code "401"
    And the value of the item "//s:message" in the response should match "/No 'Authorization: Basic' header found/"
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |


  Scenario Outline: create a new public link share of a file with edit permissions
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has uploaded file with content "Random data" to "/randomfile.txt"
    When user "Alice" creates a public link share using the sharing API with settings
      | path        | randomfile.txt            |
      | permissions | read,update,create,delete |
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And the fields of the last response to user "Alice" should include
      | item_type              | file            |
      | mimetype               | text/plain      |
      | file_target            | /randomfile.txt |
      | path                   | /randomfile.txt |
      | permissions            | read,update     |
      | share_type             | public_link     |
      | displayname_file_owner | %displayname%   |
      | displayname_owner      | %displayname%   |
      | uid_file_owner         | %username%      |
      | uid_owner              | %username%      |
      | name                   |                 |
    And the public should be able to download the last publicly shared file using the old public WebDAV API without a password and the content should be "Random data"
    And the public should be able to download the last publicly shared file using the new public WebDAV API without a password and the content should be "Random data"
    And uploading content to a public link shared file should work using the old public WebDAV API
    And uploading content to a public link shared file should work using the new public WebDAV API
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |


  Scenario Outline: creating a new public link share of a folder using the default permissions only grants read access and can be accessed with no password or any password using the public WebDAV API
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has created folder "/PARENT"
    And user "Alice" has uploaded file with content "Random data" to "/PARENT/randomfile.txt"
    When user "Alice" creates a public link share using the sharing API with settings
      | path | PARENT |
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And the fields of the last response to user "Alice" should include
      | item_type              | folder               |
      | mimetype               | httpd/unix-directory |
      | file_target            | /PARENT              |
      | path                   | /PARENT              |
      | permissions            | read                 |
      | share_type             | public_link          |
      | displayname_file_owner | %displayname%        |
      | displayname_owner      | %displayname%        |
      | uid_file_owner         | %username%           |
      | uid_owner              | %username%           |
      | name                   |                      |
    And the public should be able to download file "/randomfile.txt" from inside the last public link shared folder using the old public WebDAV API without password and the content should be "Random data"
    And the public should be able to download file "/randomfile.txt" from inside the last public link shared folder using the new public WebDAV API without password and the content should be "Random data"
    And the public upload to the last publicly shared folder using the old public WebDAV API should fail with HTTP status code "403"
    And the public upload to the last publicly shared folder using the new public WebDAV API should fail with HTTP status code "403"
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |


  Scenario Outline: creating a new public link share of a folder, with a password and accessing using the public WebDAV API
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has created folder "/PARENT"
    And user "Alice" has uploaded file with content "Random data" to "/PARENT/randomfile.txt"
    When user "Alice" creates a public link share using the sharing API with settings
      | path        | PARENT   |
      | password    | %public% |
      | permissions | change   |
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And the fields of the last response to user "Alice" should include
      | item_type              | folder               |
      | mimetype               | httpd/unix-directory |
      | file_target            | /PARENT              |
      | path                   | /PARENT              |
      | permissions            | change               |
      | share_type             | public_link          |
      | displayname_file_owner | %displayname%        |
      | displayname_owner      | %displayname%        |
      | uid_file_owner         | %username%           |
      | uid_owner              | %username%           |
      | name                   |                      |
    And the public should be able to download file "/randomfile.txt" from inside the last public link shared folder using the old public WebDAV API with password "%public%" and the content should be "Random data"
    And the public should be able to download file "/randomfile.txt" from inside the last public link shared folder using the new public WebDAV API with password "%public%" and the content should be "Random data"
    But the public should not be able to download file "/randomfile.txt" from inside the last public link shared folder using the old public WebDAV API without a password
    And the public should not be able to download file "/randomfile.txt" from inside the last public link shared folder using the old public WebDAV API with password "%regular%"
    And the public should not be able to download file "/randomfile.txt" from inside the last public link shared folder using the new public WebDAV API without a password
    And the public should not be able to download file "/randomfile.txt" from inside the last public link shared folder using the new public WebDAV API with password "%regular%"
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |

  @smokeTest
  Scenario Outline: getting the share information of public link share from the OCS API does not expose sensitive information
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has uploaded file with content "Random data" to "/randomfile.txt"
    When user "Alice" creates a public link share using the sharing API with settings
      | path     | randomfile.txt |
      | password | %public%       |
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And the fields of the last response to user "Alice" should include
      | file_target            | /randomfile.txt |
      | path                   | /randomfile.txt |
      | item_type              | file            |
      | share_type             | public_link     |
      | permissions            | read            |
      | uid_owner              | Alice           |
      | share_with             | ***redacted***  |
      | share_with_displayname | ***redacted***  |
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |


  Scenario Outline: getting the share information of password less public-links hides credential placeholders
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has uploaded file with content "Random data" to "/randomfile.txt"
    When user "Alice" creates a public link share using the sharing API with settings
      | path | randomfile.txt |
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And the fields of the last response to user "Alice" should include
      | file_target | /randomfile.txt |
      | path        | /randomfile.txt |
      | item_type   | file            |
      | share_type  | public_link     |
      | permissions | read            |
      | uid_owner   | %username%      |
    And the fields of the last response should not include
      | share_with             | ANY_VALUE |
      | share_with_displayname | ANY_VALUE |
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |


  Scenario Outline: creating a link share with no specified permissions defaults to read permissions when public upload is disabled globally and accessing using the public WebDAV API
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has created folder "/afolder"
    When user "Alice" creates a public link share using the sharing API with settings
      | path | /afolder |
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And the fields of the last response to user "Alice" should include
      | id          | A_STRING    |
      | share_type  | public_link |
      | permissions | read        |
    And the public upload to the last publicly shared folder using the old public WebDAV API should fail with HTTP status code "403"
    And the public upload to the last publicly shared folder using the new public WebDAV API should fail with HTTP status code "403"
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |


  Scenario Outline: creating a link share with edit permissions keeps it using the public WebDAV API
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has created folder "/afolder"
    When user "Alice" creates a public link share using the sharing API with settings
      | path        | /afolder                  |
      | permissions | read,update,create,delete |
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And the fields of the last response to user "Alice" should include
      | id          | A_STRING                  |
      | share_type  | public_link               |
      | permissions | read,update,create,delete |
    And uploading a file should work using the old public WebDAV API
    And uploading a file should work using the new public WebDAV API
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |


  Scenario Outline: creating a link share with upload permissions keeps it using the public WebDAV API
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has created folder "/afolder"
    When user "Alice" creates a public link share using the sharing API with settings
      | path        | /afolder    |
      | permissions | read,create |
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And the fields of the last response to user "Alice" should include
      | id          | A_STRING    |
      | share_type  | public_link |
      | permissions | read,create |
    And uploading a file should work using the old public WebDAV API
    And uploading a file should work using the new public WebDAV API
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |


  Scenario Outline: Do not allow public sharing of the root on OCIS when the default permission is read and access using the public WebDAV API
    Given using OCS API version "<ocs_api_version>"
    When user "Alice" creates a public link share using the sharing API with settings
      | path | / |
    Then the OCS status code should be "400"
    And the HTTP status code should be "<http_status_code>"
    Examples:
      | ocs_api_version | http_status_code |
      | 1               | 200              |
      | 2               | 400              |


  Scenario Outline: user creates a public link share of a file with file name longer than 64 chars using the public WebDAV API
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has uploaded file with content "long file" to "/aquickbrownfoxjumpsoveraverylazydogaquickbrownfoxjumpsoveralazydog.txt"
    When user "Alice" creates a public link share using the sharing API with settings
      | path | /aquickbrownfoxjumpsoveraverylazydogaquickbrownfoxjumpsoveralazydog.txt |
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And the public should be able to download the last publicly shared file using the old public WebDAV API without a password and the content should be "long file"
    And the public should be able to download the last publicly shared file using the new public WebDAV API without a password and the content should be "long file"

    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |


  Scenario Outline: user creates a public link share of a folder with folder name longer than 64 chars and access using the public WebDAV API
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has created folder "/aquickbrownfoxjumpsoveraverylazydogaquickbrownfoxjumpsoveralazydog"
    And user "Alice" has uploaded file with content "Random data" to "/aquickbrownfoxjumpsoveraverylazydogaquickbrownfoxjumpsoveralazydog/randomfile.txt"
    When user "Alice" creates a public link share using the sharing API with settings
      | path | /aquickbrownfoxjumpsoveraverylazydogaquickbrownfoxjumpsoveralazydog |
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And the public should be able to download file "/randomfile.txt" from inside the last public link shared folder using the old public WebDAV API without password and the content should be "Random data"
    And the public should be able to download file "/randomfile.txt" from inside the last public link shared folder using the new public WebDAV API without password and the content should be "Random data"

    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |

  @issue-1293
  Scenario: delete a folder that has been publicly shared and try to access using the public WebDAV API
    Given user "Alice" has created folder "PARENT"
    And user "Alice" has uploaded file with content "Random data" to "/PARENT/parent.txt"
    And user "Alice" has created a public link share with settings
      | path        | PARENT |
      | permissions | read   |
    When user "Alice" deletes folder "/PARENT" using the WebDAV API
    Then the public download of file "/parent.txt" from inside the last public link shared folder using the new public WebDAV API should fail with HTTP status code "404"

  @issue-1269 @issue-1293
  Scenario: try to download from a public share that has upload only permissions using the public webdav api
    Given user "Alice" has created folder "PARENT"
    And user "Alice" has uploaded file with content "Random data" to "/PARENT/parent.txt"
    And user "Alice" has created a public link share with settings
      | path        | PARENT          |
      | permissions | uploadwriteonly |
    When the public downloads file "parent.txt" from inside the last public link shared folder using the new public WebDAV API
    Then the value of the item "//s:message" in the response should be "File not found: parent.txt"
    And the HTTP status code should be "404"


  Scenario: get the size of a file shared by public link
    Given user "Alice" has uploaded file with content "This is a test file" to "test-file.txt"
    And user "Alice" has created a public link share with settings
      | path        | test-file.txt |
      | permissions | read          |
    When the public gets the size of the last shared public link using the WebDAV API
    Then the HTTP status code should be "207"
    And the size of the file should be "19"


  Scenario Outline: get the mtime of a file shared by public link
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "file.txt" with mtime "Thu, 08 Aug 2019 04:18:13 GMT" using the WebDAV API
    When user "Alice" creates a public link share using the sharing API with settings
      | path        | file.txt |
      | permissions | read     |
    Then the HTTP status code should be "200"
    And the mtime of file "file.txt" in the last shared public link using the WebDAV API should be "Thu, 08 Aug 2019 04:18:13 GMT"
    Examples:
      | dav-path-version |
      | old         |
      | new         |


  Scenario Outline: get the mtime of a file inside a folder shared by public link
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "testFolder"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "testFolder/file.txt" with mtime "Thu, 08 Aug 2019 04:18:13 GMT" using the WebDAV API
    When user "Alice" creates a public link share using the sharing API with settings
      | path        | /testFolder |
      | permissions | read        |
    Then the HTTP status code should be "200"
    And the mtime of file "file.txt" in the last shared public link using the WebDAV API should be "Thu, 08 Aug 2019 04:18:13 GMT"
    Examples:
      | dav-path-version |
      | old         |
      | new         |


  Scenario: get the mtime of a file inside a folder shared by public link using new webDAV version
    Given user "Alice" has created folder "testFolder"
    And user "Alice" has created a public link share with settings
      | path        | /testFolder               |
      | permissions | read,update,create,delete |
    When the public uploads file "file.txt" to the last public link shared folder with mtime "Thu, 08 Aug 2019 04:18:13 GMT" using the new public WebDAV API
    Then the HTTP status code should be "201"
    And as "Alice" file "testFolder/file.txt" should exist
    And as "Alice" the mtime of the file "testFolder/file.txt" should be "Thu, 08 Aug 2019 04:18:13 GMT"
    And the mtime of file "file.txt" in the last shared public link using the WebDAV API should be "Thu, 08 Aug 2019 04:18:13 GMT"


  Scenario: overwriting a file changes its mtime (new public webDAV API)
    Given user "Alice" has created folder "testFolder"
    When user "Alice" uploads file with content "uploaded content for file name ending with a dot" to "testFolder/file.txt" using the WebDAV API
    And user "Alice" creates a public link share using the sharing API with settings
      | path        | /testFolder               |
      | permissions | read,update,create,delete |
    And the public uploads file "file.txt" to the last public link shared folder with mtime "Thu, 08 Aug 2019 04:18:13 GMT" using the new public WebDAV API
    Then the HTTP status code should be "204"
    And as "Alice" file "/testFolder/file.txt" should exist
    And as "Alice" the mtime of the file "testFolder/file.txt" should be "Thu, 08 Aug 2019 04:18:13 GMT"
    And the mtime of file "file.txt" in the last shared public link using the WebDAV API should be "Thu, 08 Aug 2019 04:18:13 GMT"


  Scenario: check the href of a public link file
    Given using new DAV path
    And user "Alice" has uploaded file with content "Random data" to "/file.txt"
    And user "Alice" has created a public link share with settings
      | path        | file.txt |
      | permissions | read     |
    When the public lists the resources in the last created public link with depth "1" using the WebDAV API
    Then the HTTP status code should be "207"
    And the value of the item "//d:response[2]//d:href" in the response should match "/%base_path%\/remote.php\/dav\/public-files\/%public_token%\/file.txt$/"
    When the public gets the following properties of entry "/file.txt" in the last created public link using the WebDAV API
      | propertyName |
      | d:href       |
    Then the HTTP status code should be "207"
    And the value of the item "//d:href" in the response should match "/%base_path%\/remote.php\/dav\/public-files\/%public_token%\/file.txt$/"
