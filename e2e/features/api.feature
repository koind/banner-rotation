Feature: microservices for banner-rotation implementation
  There should be an opportunity to add, remove a banner from rotation.
  Also possible to set the transition for banners in rotation and choose a banner to display.
  The should also be sending statistics to the queue.

  Scenario: Add banner1 to rotation
    When 1. I send "POST" request to "http://api:7766/banner/add" with "application/json" data:
    """
    {
        "bannerId": 1,
        "slotId": 1,
        "description": "banner 1"
    }
    """
    Then The response code should be 200

  Scenario: Select banner to display
    When 2. I send "POST" request to "http://api:7766/banner/select" with "application/json" data:
    """
    {
        "slotId": 1,
        "groupId": 1
    }
    """
    Then The response code should be 200
    And The response should match id banner "1"


  Scenario: Set transition for banner
    When 3. I send "POST" request to "http://api:7766/banner/set-transition" with "application/json" data:
    """
    {
        "bannerId": 1,
        "groupId": 1
    }
    """
    Then The response code should be 200

  Scenario: Add banner2 to rotation
    When 4. I send "POST" request to "http://api:7766/banner/add" with "application/json" data:
    """
    {
        "bannerId": 2,
        "slotId": 1,
        "description": "banner 2"
    }
    """
    Then The response code should be 200

  Scenario: Select banner to display
    When 5. I send "POST" request to "http://api:7766/banner/select" with "application/json" data:
    """
    {
        "slotId": 1,
        "groupId": 1
    }
    """
    Then The response code should be 200
    And The response should match id banner "2"

  Scenario: Select banner to display
    When 6. I send "POST" request to "http://api:7766/banner/select" with "application/json" data:
    """
    {
        "slotId": 1,
        "groupId": 1
    }
    """
    Then The response code should be 200
    And The response should match id banner "1"

  Scenario: Removes the banner from the rotation
    When 7. I send "DELETE" request to "http://api:7766/banner/remove/1"
    Then The response code should be 200