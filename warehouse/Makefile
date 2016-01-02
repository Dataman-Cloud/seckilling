default: help

py-cleanup:
	-find . -name '*.pyc' -delete
	-find . -name '*.pyo' -delete
	-find . -name '__pycache__' -delete

init:
	docker-compose -p omega_god up -d

# TODO Improve me by docker-compose one-off cmd
create-db:
	docker exec -it omegagod_omega_god_db_1 mysql -uroot -p111111 -e "create database omega_god"
	mysql -uroot -p111111 -h127.0.0.1 -e "create database omega"

# TODO Improve me by docker-compose one-off cmd
init-db:
	docker exec -it omegagod_omega_god_1 python3 manage.py migrate
	mysql -uroot -p111111 -h127.0.0.1 omega < mysql_settings/omega_dump

# TODO Improve me by docker-compose one-off cmd
create-superuser:
	docker exec -it omegagod_omega_god_1 python3 manage.py createsuperuser

restart-god:
	docker restart omegagod_omega_god_1

cleanup: py-cleanup
	docker-compose -p omega_god stop
	docker-compose -p omega_god rm -f

help:
	@echo 'Usage make [command]'
	@echo ''
	@echo 'Available commands:'
	@echo ''
	@echo '  init                    - start-containers and update-ip-port'
	@echo '  cleanup                 - Clean up development env gracefully'
	@echo '  create-db               - Create db'
	@echo '  create-superuser        - Create superuser'
	@echo '  generate-revision-for-db - auto generate revision for db.'

	@echo '  py-cleanup              - Clean up *.pyc, *.pyo, __pycache__'
	@echo '  help                    - Show this help message and exit'

