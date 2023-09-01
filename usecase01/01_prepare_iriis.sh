source ./path.sh $1
RAW_DIR="$IRIIS_DIR/raw"
REG_DIR="$IRIIS_DIR/regionals"

for reg in $(seq 12); do
    > $REG_DIR/REG$reg.csv
    find $RAW_DIR/*REGIONAL${reg}_* -type f 2> /dev/null 1> /dev/null \
    &&  for file in $(find $RAW_DIR/*REGIONAL${reg}_* -type f 2> /dev/null); do
            paste -d, <(cut -d, -f8,9 $file \
                            | tr -s , -) <(cut -d, -f1,4 $file) \
                | uniq \
                | tail +2 \
                >> $REG_DIR/REG$reg.csv
        done \
    ||  rm $REG_DIR/REG$reg.csv
done