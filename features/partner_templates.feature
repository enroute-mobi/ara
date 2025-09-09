Feature: Manages Partner templates

  Background:
    Given a Referential "test" is created

  Scenario: Create a Partner template
    Given a Partner template "test" exists with the following attributes:
      | CredentialType   | FormatMatching                                                                   |
      | LocalCredential  | %{value}                                                                         |
      | RemoteCredential | %{value}                                                                         |
      | ConnectorTypes   | ["siri-check-status-client","siri-estimated-timetable-subscription-broadcaster"] |
      | Settings         | {"remote_code_space":"internal"}                                                 |
    Then one Partner template has the following attributes:
      | Slug | test |

  Scenario: Update a Partner template
    Given a Partner template "test" exists with the following attributes:
      | CredentialType   | FormatMatching                                                                   |
      | LocalCredential  | %{value}                                                                         |
      | RemoteCredential | %{value}                                                                         |
      | ConnectorTypes   | ["siri-check-status-client","siri-estimated-timetable-subscription-broadcaster"] |
      | Settings         | {"remote_code_space":"internal"}                                                 |
    When the Partner template "test" is updated with the following attributes:
      | LocalCredential | test:%{value} |
    Then one Partner template has the following attributes:
      | Slug            | test          |
      | LocalCredential | test:%{value} |
