package targets

type Target interface {
	isTarget() // данный метод нужен как заглушка, чтобы вы в функцию не передали любой рандомный тип, а только типы StaticTarget или DynamicTarget
}

type StaticTarget struct {
	Host string
}

func (s StaticTarget) isTarget() {}

type DynamicTarget struct {
	Cluster   string
	Namespace string
	Host      string
}

func (d DynamicTarget) isTarget() {}

// Эта функция нужна для передачи слайсов интерфейсных типов в фунцию GetMetadata
func ToTargets[T Target](targets []T) []Target {
	var targetsSlice []Target = make([]Target, len(targets))
	for i, target := range targets {
		targetsSlice[i] = target
	}
	return targetsSlice
}
