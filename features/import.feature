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
