DATE=$1
DATA_PATH="/home/agung/storage/local-data/tsel-medserver/dataset-03"
SCRIPT_PATH="/home/agung/storage/devs/tsel-medserver/usecase02"
WORKSPACE_PATH=$(echo $SCRIPT_PATH/workspaces/DATE | sed "s/DATE/$DATE/")

# constants
RESULT_COLS="$SCRIPT_PATH/static/result-cols.txt"
WORKSPACE_RETENTION=7

# raw files
IRONMAN_FILENAME=$(echo $DATA_PATH/DATE_newsite.csv | sed "s/DATE/$DATE/")
CROSSFEEDER_DIR="$DATA_PATH/crossfeeder"
CROSSFEEDER_FILENAME="cross-feeder-DATE.csv"

# intermediate dirs
CROSSFEEDER_CELLS="$WORKSPACE_PATH/01_crossfeeder_cells"
IRONMAN_CELLS="$WORKSPACE_PATH/02_ironman_cells"
JOIN_IRONMAN_CROSSFEEDER_DIR="$WORKSPACE_PATH/03_result"
RESULTS="$SCRIPT_PATH/results"