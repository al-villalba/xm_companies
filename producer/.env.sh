set -o allexport

if [ ! -f "./.env" ]; then
	echo "ERROR: .env not found!"
	exit 1
fi
source .env

if [ -d .git ]; then
	GIT_COMMIT=$(git rev-parse HEAD)
else
	GIT_COMMIT=""
fi

set +o allexport

#go_run() {
#	$(find . -maxdepth 1 -perm -u=x -type f | head -n 1)
#}
