
# TODO проверять на установку генератора
# Генерирует хэндлеры на основе openapi спецификации
server-api-gen:
	oapi-codegen --config=server/openapi-gen-cfg.yaml openapi.yaml