DATE=$1
REGIONAL=$2
DATA_PATH="/home/agung/storage/local-data/tsel-medserver/dataset-03"
SCRIPT_PATH="/home/agung/storage/devs/tsel-medserver/usecase01"
WORKSPACE_PATH=$(echo $SCRIPT_PATH/workspaces/DATE | sed "s/DATE/$DATE/")

# constants
RSRPMIN=-105.0
RATIOMIN=0.8
WORKSPACE_RETENTION=7
RESULT_COLS="$SCRIPT_PATH/static/result-cols.txt"
RESULT_BADGRIDS_COLS="$SCRIPT_PATH/static/result-badgrids-cols.txt"
RESULT_IRIIS_UNDEFINED_CELLS_COLS="$SCRIPT_PATH/static/result-iriis-undefined-cells-cols.txt"

# raw files
IRONMAN_FILENAME=$(echo $DATA_PATH/DATE_newsite.csv | sed "s/DATE/$DATE/")
IRIIS_DIR="$DATA_PATH/iriis"
IRIIS_RAW_DIR="$IRIIS_DIR/raw"
IRIIS_REGIONALS_DIR="$IRIIS_DIR/regionals"
GEOHASH_RAW_DIR="$DATA_PATH/geohash-20220519"
GEOHASH_RAW_FILENAME="geohash7_NoName_Geolocation_Daily_4G_RREGION_DATE.csv"

# intermediate dirs
IRONMAN_REGIONALS="$WORKSPACE_PATH/02_ironman_regionals"
IRIIS_CELLS="$WORKSPACE_PATH/03_iriis_cells"
GEOHASH_DIR="$WORKSPACE_PATH/04_geohash_files"
JOIN_IRIIS_GEOHASH_DIR="$WORKSPACE_PATH/05_join_iriis_mediation"
CALCULATION_RESULT_DIR="$WORKSPACE_PATH/06_calculation_results"
COMBINED_DAYS_DIR="$WORKSPACE_PATH/07_combined_days"
COMBINED_REGIONS_DIR="$WORKSPACE_PATH/08_combined_regions"
RESULTS="$SCRIPT_PATH/results"