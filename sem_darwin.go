package semaphore

import "errors"

func (s *Semaphore) GetValue() (int, error){
	return 0, errors.New("not implement")
}
