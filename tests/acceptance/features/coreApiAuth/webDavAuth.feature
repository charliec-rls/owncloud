@api
Feature: auth
  As a user
  I want to check the authentication of the application
  So that I can make sure it's secure

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files

  @smokeTest
  Scenario Outline: using WebDAV anonymously
    When a user requests "<dav_path>" with "PROPFIND" and no authentication
    Then the HTTP status code should be "401"
    Examples:
      | dav_path              |
      | /remote.php/webdav    |
      | /dav/spaces/%spaceid% |

  @smokeTest
  Scenario Outline: using WebDAV with basic auth
    When user "Alice" requests "<dav_path>" with "PROPFIND" using basic auth
    Then the HTTP status code should be "207"
    Examples:
      | dav_path           |
      | /remote.php/webdav |

    @skipOnRevaMaster
    Examples:
      | dav_path              |
      | /dav/spaces/%spaceid% |
