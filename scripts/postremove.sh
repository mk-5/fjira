rm -rf "$HOME/.fjira"
for dir in /home/*/.fjira
do
  if [ -d "$dir" ]; then
    echo "Removing $dir"
    rm -rf $dir
  fi;
done
