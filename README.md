# Simple Go Data Generator

This is a simple command-line tool to generate large amounts of highly randomized data.
Just run `datagen -f yamlconfig.yml`

Here is a sample YAML configuration file:

```
directory:         /usr/local/dataset/pumpdata
zipFileName:       pumpdata.zip
numberOfEntities:  2
totalTimeInHours:  6
pkFileName:        water_pumps.csv
files:
   - fileName:       pressure_psi.csv
     dataType:       timeseries
     timeStepMillis: 1500
     minValue:       1000
     maxValue:       1300

   - fileName:       manufacturer.csv
     dataType:       static
     values:         [ Flowserve, Yildiz, Andoria, Enerpac, SNC, Condor, Hankia, Delta ]

   - fileName:       vibration_khz.csv
     dataType:       timeseries
     timeStepMillis: 60000
     minValue:       150
     maxValue:       230

   - fileName:       water_flow_liters_per_minute.csv
     dataType:       timeseries
     timeStepMillis: 360000
     minValue:       9.2
     maxValue:       10.5

   - fileName:       system_lubrication_level_percent.csv
     dataType:       timeseries
     timeStepMillis: 3600000
     minValue:       0.5
     maxValue:       1.0

   - fileName:       location.csv
     dataType:       static
     values:         [ Warehouse Floor, Store Room, Warehouse Basement, Outside, Shed ]

   - fileName:       relative_humidity.csv
     dataType:       timeseries
     timeStepMillis: 21600000
     minValue:       0.2
     maxValue:       0.5
```
