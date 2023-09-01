source ./path.sh $1
DATE=$1

mkdir -p $IRONMAN_REGIONALS
REGS=$(tail +2 $IRONMAN_FILENAME | cut -d'|' -f1 | sort -n -u)
> $IRONMAN_REGIONALS/index.txt
for reg in $REGS; do
    echo $reg >> $IRONMAN_REGIONALS/index.txt
    cat $IRONMAN_FILENAME \
        | tr -s '|' , \
        | grep ^$reg, \
        | sort -u \
        | cut -d, -f2,3 \
        | tr -s , - \
        > $IRONMAN_REGIONALS/cells_R$reg.csv
    echo "  regional $reg: "$(cat $IRONMAN_REGIONALS/cells_R$reg.csv | wc -l)
done