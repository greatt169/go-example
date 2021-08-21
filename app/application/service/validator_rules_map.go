package service

// Правила валидации входных данных
type ValidatorRulesMap struct{}

// Возвращает правила валидации для запроса списка новостей
func (v *ValidatorRulesMap) GetGetNewsMap() map[string][]string {
	mapRules := map[string][]string{
		"limit":  {"numeric", "min:0", "max:100"},
		"offset": {"numeric"},
		"sort":   {"in:active_from,date_create"},
		"order":  {"in:asc,desc"},
		"query":  {"min_string_len:3", "max_string_len:100"},
		"mode":   {"in:active,inactive"},
	}
	return mapRules
}

// Возвращает правила валидации для запроса создания новости
func (v *ValidatorRulesMap) GetCreateNewsMap() map[string][]string {
	mapRules := map[string][]string{
		"title":       {"required", "max_string_len:200"},
		"text":        {"required"},
		"textJson":    {"required"},
		"activeFrom":  {"required", "numeric"},
		"isImportant": {"bool"},
	}
	return mapRules
}

// Возвращает правила валидации для создания тега новости
func (v *ValidatorRulesMap) GetCreateNewsTagMap() map[string][]string {

	mapRules := map[string][]string{
		"name": {"required", "max_string_len:50", "regex:[a-zа-я0-9]+$"},
	}
	return mapRules
}

// Возвращает правила валидации для обновления тега новости
func (v *ValidatorRulesMap) GetUpdateNewsTagMap() map[string][]string {
	return v.GetCreateNewsTagMap()
}

//  Возвращает список правил для проверки файлов новостей
func (v *ValidatorRulesMap) GetNewsFilesMap() map[string]string {
	mapRules := map[string]string{
		"file_ext":    "doc,xls,ppt,jpg,bmp,pdf,rtf,txt,zip",
		"file_size":   "5242880",
		"files_limit": "7",
	}
	return mapRules
}

// Возвращает правила валидации для запроса обновления новости
func (v *ValidatorRulesMap) GetUpdateNewsMap() map[string][]string {
	return v.GetCreateNewsMap()
}

// Возвращает правила валидации для запроса списка акций
func (v *ValidatorRulesMap) GetGetPromoMap() map[string][]string {
	mapRules := map[string][]string{
		"limit":  {"numeric", "min:0", "max:100"},
		"sort":   {"in:active_from,date_create"},
		"order":  {"in:asc,desc"},
		"offset": {"numeric"},
	}
	return mapRules
}

// Возвращает правила валидации для запроса создания акции
func (v *ValidatorRulesMap) GetCreatePromoMap() map[string][]string {
	mapRules := map[string][]string{
		"title":      {"required", "max_string_len:140"},
		"text":       {"required"},
		"textJson":   {"required"},
		"activeFrom": {"required", "numeric"},
		"period":     {"max_string_len:80"},
	}
	return mapRules
}

// Возвращает правила валидации для запроса обновления акции
func (v *ValidatorRulesMap) GetUpdatePromoMap() map[string][]string {
	return v.GetCreatePromoMap()
}

// Возвращает правила валидации для запроса получения детальной элемента по uuid
func (v *ValidatorRulesMap) GetGetOneMap() map[string][]string {
	mapRules := map[string][]string{
		"id": {"uuid_v4"},
	}
	return mapRules
}

// Возвращает правила валидации для запроса получения детальной элемента по slug
func (v *ValidatorRulesMap) GetGetOneBySlugMap() map[string][]string {
	mapRules := map[string][]string{
		"slug": {"required"},
	}
	return mapRules
}

// Возвращает правила валидации для запроса удаления элемента по uuid
func (v *ValidatorRulesMap) GetDeleteMap() map[string][]string {
	return v.GetGetOneMap()
}

// Возвращает правила валидации для запроса подписки
func (v *ValidatorRulesMap) GetSubscribeMap() map[string][]string {
	mapRules := map[string][]string{
		"email": {"email"},
	}
	return mapRules
}
