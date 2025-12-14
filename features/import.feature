Feature: Import with Referential API

  @ARA-890
  Scenario: Import by using the API token
    Given a Referential "test" exists with the following attributes:
      | Tokens | dummy,another |
    When I import in the referential "test" with the token "dummy" these models:
      """
stop_area,5381c0d7-a479-4f6e-a5e8-36072200715c,"","",2031-01-04,First Stop Place,"{""external"":""ABC"",""internal"":""123""}",[],{},{},true,false,true
      """
    Then the import should be successful

  @ARA-890
  Scenario: Import by using the import token
    Given a Referential "test" exists with the following attributes:
      | Import Tokens | dummy,another |
    When I import in the referential "test" with the token "dummy" these models:
      """
stop_area,5381c0d7-a479-4f6e-a5e8-36072200715c,"","",2031-01-04,First Stop Place,"{""external"":""ABC"",""internal"":""123""}",[],{},{},true,false,true
      """
    Then the import should be successful

  @ARA-890
  Scenario: Import by using a wrong token
    Given a Referential "test" exists with the following attributes:
      | Import Tokens | another |
    When I import in the referential "test" with the token "dummy" these models:
      """
stop_area,5381c0d7-a479-4f6e-a5e8-36072200715c,"","",2031-01-04,First Stop Place,"{""external"":""ABC"",""internal"":""123""}",[],{},{},true,false,true
      """
    Then the import must fail with an unauthorized status

  Scenario: Import without aimed stop visit counts
    Given a Referential "test" exists with the following attributes:
      | Import Tokens | dummy |
    When I import in the referential "test" with the token "dummy" these models:
      """
line,bfbfa8d5-9988-4452-9771-28cf6ab3706a,2031-01-01,First Line,{},{},{},true
vehicle_journey,01eebc99-9c0b-4ef8-bb6d-6bb9bd380a11,2031-01-01,Name,{},bfbfa8d5-9988-4452-9771-28cf6ab3706a,origin,destination,{},{},outbound
      """
    Then the import should be successful

  Scenario: Import with aimed stop visit counts
    Given a Referential "test" exists with the following attributes:
      | Import Tokens | dummy |
    When I import in the referential "test" with the token "dummy" these models:
      """
line,bfbfa8d5-9988-4452-9771-28cf6ab3706a,2031-01-01,First Line,{},{},{},true
vehicle_journey,01eebc99-9c0b-4ef8-bb6d-6bb9bd380a11,2031-01-01,Name,{},bfbfa8d5-9988-4452-9771-28cf6ab3706a,origin,destination,{},{},outbound,42
      """
    Then the import should be successful
