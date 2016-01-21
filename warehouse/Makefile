default: help

py-cleanup:
	-find . -name '*.pyc' -delete
	-find . -name '*.pyo' -delete
	-find . -name '__pycache__' -delete

init:
	docker-compose -p warehouse up -d

# TODO Improve me by docker-compose one-off cmd
init-db:
	docker exec -it warehouse_app_1 python3 manage.py migrate

# TODO Improve me by docker-compose one-off cmd
create-superuser:
	docker exec -it warehouse_app_1 python3 manage.py createsuperuser

restart:
	docker restart warehouse_app_1
	docker restart warehouse_nginx_1

cleanup: py-cleanup
	docker-compose -p warehouse stop
	docker-compose -p warehouse rm -f

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
