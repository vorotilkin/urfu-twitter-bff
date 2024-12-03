
# TODO проверять на установку генератора
# Генерирует хэндлеры на основе openapi спецификации
openapi-gen:
	oapi-codegen --config=openapigen/openapi-gen-cfg.yaml openapi.yaml