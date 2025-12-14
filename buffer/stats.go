package buffer

type PoolStats struct {
	// TotalAlloc — суммарное количество байт,
	// аллоцированных при использовании пула.
	TotalAlloc int64

	// TotalCrops — суммарное количество байт,
	// обрезанных при возвратах в пул буферов
	// большего размера, чем требуется.
	TotalCrops int64

	// TotalOverflows — суммарное количество
	// попыток возврата буфера в заполненный пул.
	TotalOverflows int64

	// InUse — суммарный размер всех буферов в пуле
	InUse int64
}
