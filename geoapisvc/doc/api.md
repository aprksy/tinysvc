# API Enpoints
## Get sites of provided cells
  |path|method|type|content-type|
  |---|---|---|---|
  `/cells/site`|POST|atomic|application/json|

  ### example
  #### request body:
  ```
  {
    "date": "20220520",
    "region": "01",
    "cells": "LGS092MT1_PMKIMANPKAMLAMALANGSAMT03,LGS101ML1"
  }
  ```
  #### response body:
  ```
  {
    "LGS092MT1_PMKIMANPKAMLAMALANGSAMT03": "LGS092",
    "LGS101ML1": "LGS101"
  }
  ```
## Get details of provided cells
  |path|method|type|content-type|
  |---|---|---|---|
  `/cells/details`|POST|atomic|application/json|

  ### example
  #### request body:
  ```
  {
    "date": "20220520",
    "region": "01",
    "cells": "LGS092MT1_PMKIMANPKAMLAMALANGSAMT03,LGS101ML1"
  }
  ```
  #### response body:
  ```
  {
    "LGS092MT1_PMKIMANPKAMLAMALANGSAMT03": {
        "azimuth": "255",
        "band": "900",
        "cell": "LGS092MT1_PMKIMANPKAMLAMALANGSAMT03",
        "et": "6",
        "hbw": "35",
        "height": "24",
        "id": 1,
        "model": "Mobi_MB4BQMF651616DEINTSL",
        "mt": "2",
        "tech": "4G",
        "vendor": "ERICSSON"
    },
    "LGS101ML1": {
        "azimuth": "100",
        "band": "1800",
        "cell": "LGS101ML1",
        "et": "2",
        "hbw": "40",
        "height": "70",
        "id": 2,
        "model": "Mobi_MB4BQMF651616DEINTSL",
        "mt": "0",
        "tech": "4G",
        "vendor": "ERICSSON"
    }
  }
  ```
## Get cells of provided sites
  |path|method|type|content-type|
  |---|---|---|---|
  `/sites/cells`|POST|atomic|application/json|

  ### example
  #### request body:
  ```
  {
    "date": "20220520",
    "region": "01",
    "sites": "LGS092,LGS101"
  }
  ```
  #### response body:
  ```
  {
    "LGS092": [
        "LGS092ML1",
        "LGS092MT1_PMKIMANPKAMLAMALANGSAMT01",
        "LGS092ML2",
        "LGS092MT1_PMKIMANPKAMLAMALANGSAMT02",
        "LGS092ML3",
        "LGS092MT1_PMKIMANPKAMLAMALANGSAMT03",
        "LGS092MR1",
        "LGS092MR2",
        "LGS092MR3"
    ],
    "LGS101": [
        "LGS101ML1",
        "LGS101MT1",
        "LGS101ML2",
        "LGS101MT2",
        "LGS101ML3",
        "LGS101MT3",
        "LGS101MR1",
        "LGS101ME1",
        "LGS101MF1",
        "LGS101MR2",
        "LGS101ME2",
        "LGS101MF2",
        "LGS101MR3"
    ]
  }
  ```
## Get details of provided sites
  |path|method|type|content-type|
  |---|---|---|---|
  `/sites/details`|POST|atomic|application/json|

  ### example
  #### request body:
  ```
  {
    "date": "20220520",
    "region": "01",
    "sites": "LGS092,LGS101"
  }
  ```
  #### response body:
  ```
  {
    "LGS092": {
        "id": 1,
        "lat": "4.461801",
        "lon": "97.967813",
        "name": "LGS092MM2-PMKIMAN-P-KAM-LAMA-LANGSA",
        "site": "LGS092",
        "type": "Macro"
    },
    "LGS101": {
        "id": 2,
        "lat": "4.506755",
        "lon": "97.877305",
        "name": "LGS101MM2-ALUETEH",
        "site": "LGS101",
        "type": "Macro"
    }
  }
  ```
## Get sites that intersect provided boundary
  |path|method|type|content-type|
  |---|---|---|---|
  `/sites/intersects`|POST|atomic|application/json|

  ### example
  #### request body:
  ```
  {
    "date": "20220520",
    "region": "01",
    "areaType": "circle",
    "data": {
        "lat": 3.580112249763098,
        "lng": 98.69898319244385,
        "radius": 1000
    }
  }
  ```
  ```
  {
    "date": "20220520",
    "region": "01",
    "areaType": "geojson",
    "data": {
        "type": "FeatureCollection",
        "features": [
        {
            "type": "Feature",
            "properties": {},
            "geometry": {
            "type": "Point",
            "coordinates": [
                98.69898319244385,
                3.580112249763098
            ]
            }
        },
        {
            "type": "Feature",
            "properties": {},
            "geometry": {
            "type": "Polygon",
            "coordinates": [
                [
                [
                    98.6962366104126,
                    3.5848237129085265
                ],
                [
                    98.69404792785645,
                    3.580026586572354
                ],
                [
                    98.69563579559326,
                    3.5757862486190333
                ],
                [
                    98.70108604431152,
                    3.575871912206189
                ],
                [
                    98.70323181152344,
                    3.5841812420899286
                ],
                [
                    98.6962366104126,
                    3.5848237129085265
                ]
                ]
            ]
            }
        }
        ]
    }
  }
  ```
  #### response body:
  ```
  [
    {
        "Id": "MDN504",
        "Lat": 3.604166,
        "Lng": 98.690833,
        "Fields": {}
    },
    {
        "Id": "MDN520",
        "Lat": 3.5913880000000002,
        "Lng": 98.696666,
        "Fields": {}
    },
    {
        "Id": "MDN521",
        "Lat": 3.531971,
        "Lng": 98.700961,
        "Fields": {}
    }
    ...
  ]
  ```
## Get cell and tile of provided event ids
  |path|method|type|content-type|
  |---|---|---|---|
  `/events/cell-tile`|POST|atomic|application/json|

  ### example
  #### request body:
  ```
  {
    "date": "20220520",
    "region": "01",
    "events": "1961135612710947840,1962107426371032576,1962107426371041280"
  }
  ```
  #### response body:
  ```
  {
    "1961135612710947840": {
        "cell": "E_MDX226HL1_MCJlSutrisno-TBG_HL02",
        "id": 1,
        "tile": "5764608485849889437"
    },
    "1962107426371032576": {
        "cell": "E_MDX226HL1_MCJlSutrisno-TBG_HL03",
        "id": 2,
        "tile": "5764608485849889341"
    },
    "1962107426371041280": {
        "cell": "E_MDX226HL1_MCJlSutrisno-TBG_HL03",
        "id": 3,
        "tile": "5764608485849889341"
    }
  }
  ```
## Get events that intersect provided boundary
  |path|method|type|content-type|
  |---|---|---|---|
  `/events/intersects`|POST|atomic|application/json|

  ### example
  #### request body:
  ```
  {
    "date": "20220520", 
    "region": "01", 
    "areaType": "circle", 
    "data": {
        "lat": 3.580112249763098, 
        "lng": 98.69898319244385, 
        "radius": 50
    }
  }
  ```
  ```
  {
    "date": "20220520",
    "region": "01",
    "areaType": "geojson",
    "data": {
        "type": "FeatureCollection",
        "features": [
        {
            "type": "Feature",
            "properties": {},
            "geometry": {
            "type": "Point",
            "coordinates": [
                98.69898319244385,
                3.580112249763098
            ]
            }
        },
        {
            "type": "Feature",
            "properties": {},
            "geometry": {
            "type": "Polygon",
            "coordinates": [
                [
                [
                    98.6962366104126,
                    3.5848237129085265
                ],
                [
                    98.69404792785645,
                    3.580026586572354
                ],
                [
                    98.69563579559326,
                    3.5757862486190333
                ],
                [
                    98.70108604431152,
                    3.575871912206189
                ],
                [
                    98.70323181152344,
                    3.5841812420899286
                ],
                [
                    98.6962366104126,
                    3.5848237129085265
                ]
                ]
            ]
            }
        }
        ]
    }
  }
  ```
  #### response body:
  ```
  [
    {
        "Id": "1951385757778559488",
        "Lat": 98.698929,
        "Lng": 3.580205,
        "Fields": {
        "latt": 3.580069,
        "long": 98.69894,
        "rsrp": -110
        }
    },
    {
        "Id": "1961735752788489728",
        "Lat": 98.698929,
        "Lng": 3.580205,
        "Fields": {
        "latt": 3.580069,
        "long": 98.69894,
        "rsrp": -113
        }
    },
    {
        "Id": "1969667938907471872",
        "Lat": 98.698929,
        "Lng": 3.580205,
        "Fields": {
        "latt": 3.580069,
        "long": 98.69894,
        "rsrp": -101
        }
    },
    ...
  ]
  ```
## Add boundary
  |path|method|type|content-type|
  |---|---|---|---|
  `/boundary`|PUT|atomic|application/json|

  ### example
  #### request body:
  ```
  {
      "id": "circle-2", 
      "areaType": "circle", 
      "data": {
          "lat": 3.580112249763098, 
          "lng": 98.69898319244385, 
          "radius": 100
      }
  }
  ```
  #### response body: `none`

  #### request body:
  ```
  {
    "id": "polygon-2",
    "areaType": "geojson",
    "data": {
        "type": "FeatureCollection",
        "features": [
        {
            "type": "Feature",
            "properties": {},
            "geometry": {
            "type": "Point",
            "coordinates": [
                98.69898319244385,
                3.580112249763098
            ]
            }
        },
        {
            "type": "Feature",
            "properties": {},
            "geometry": {
            "type": "Polygon",
            "coordinates": [
                [
                [
                    98.6962366104126,
                    3.5848237129085265
                ],
                [
                    98.69404792785645,
                    3.580026586572354
                ],
                [
                    98.69563579559326,
                    3.5757862486190333
                ],
                [
                    98.70108604431152,
                    3.575871912206189
                ],
                [
                    98.70323181152344,
                    3.5841812420899286
                ],
                [
                    98.6962366104126,
                    3.5848237129085265
                ]
                ]
            ]
            }
        }
        ]
    }
  }
  ```
  #### response body: `none`

## Get boundary
  |path|method|type|content-type|
  |---|---|---|---|
  `/boundary?id=circle-1`|GET|atomic|-|

  ### example
  #### request body: `none`
  #### response body:
  ```
  {
      "id": "circle-1", 
      "areaType": "circle", 
      "data": {
          "lat": 3.580112249763098, 
          "lng": 98.69898319244385, 
          "radius": 100
      }
  }
  ```

## Update boundary
  |path|method|type|content-type|
  |---|---|---|---|
  `/boundary`|GET|atomic|application/json|

  ### example
  #### request body:
  ```
  {
      "id": "circle-1", 
      "areaType": "circle", 
      "data": {
          "lat": 3.580112249763098, 
          "lng": 98.69898319244385, 
          "radius": 100
      }
  }
  ```
  #### response body: `none`

## Delete boundary
  |path|method|type|content-type|
  |---|---|---|---|
  `/boundary`|GET|atomic|application/json|

  ### example
  #### request body:
  ```
  {
      "id": "circle-2"
  }
  ```
  #### response body: `none`

## Get Boundary Facts
  |path|method|type|content-type|
  |---|---|---|---|
  `/bfact`|POST|atomic|application/json|

  ### example
  #### request body:
  ```
  {
    "date": "20220520",
    "region": "01",
    "boundaryId": "polygon-1"
  }
  ```
  #### request body:
  ```
  {
    "cacheId": "OAzBA3YU4JKl0mqK"
  }
  ```
  #### response body:
  ```
  {
    "events": [
        {
            "Id": "1965955064264557568",
            "Lat": 3.576224,
            "Lng": 98.698789,
            "Fields": {
                "latt": 3.5763,
                "long": 98.69894,
                "rsrp": -111
            }
        },
        ...
    ],
    "resultId": "hWQk2rRqXO5R9G5v",
    "sites": [
        {
            "Id": "MDX226",
            "Lat": 3.5818849999999998,
            "Lng": 98.698167,
            "Fields": {}
        },
        ...
    ]
  }
  ```