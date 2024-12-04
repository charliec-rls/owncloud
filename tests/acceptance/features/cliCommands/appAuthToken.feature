@env-config
Feature: auth-app token


  Background:
    Given the following configs have been set:
      | config                | value    |
      | OCIS_ADD_RUN_SERVICES | auth-app |
      | PROXY_ENABLE_APP_AUTH | true     |
    And user "Alice" has been created with default attributes and without skeleton files

  @issue-9498
  Scenario: check backup consistency after uploading a file multiple times
    When the administrator checks the backup consistency using the CLI
    And the administrator creates app token for user "Alice" with expiration time "72h" using the CLI
    Then the command should be successful
    And the command output should be "App token created for Alice\r\n\stoken: %created_app_token%"