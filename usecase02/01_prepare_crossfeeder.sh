source ./path.sh $1
DATE=$1 
DATE1=$2

mkdir -p $CROSSFEEDER_CELLS
FNAME=$CROSSFEEDER_DIR/$(echo $CROSSFEEDER_FILENAME | sed "s/DATE/$DATE1/1")
tail +2 $FNAME \
    | sed 's/\"[0-9]*@[0-9]*,[0-9]*@[0-9]*\"/XXX/' \
    | cut -d, -f4,6,21 \
    | sed 's/,/-/1' \
    | sort -u \
    > $CROSSFEEDER_CELLS/cells-$DATE1.csv