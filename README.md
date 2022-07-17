# Flight Tracker
  API to process array of flights and generate routes for a given person

## Description
  API Method
    ``POST /track``


  API Request Params (JSON Format)
     * source `string`
     * flights `array of strings (source and destination)`

     *Example
        ```
          {
              "source": "SFO",
              "flights": [
                  ["IND", "EWR"], ["SFO", "ATL"], ["GSO", "IND"], ["ATL", "GSO"]
              ]
          }
        ```

  API Response (JSON Format)
     * airport quote `array of strings for stops`
        ```
          [
            "SFO",
            "ATL",
            "GSO",
            "IND",
            "EWR"
          ]
        ```


