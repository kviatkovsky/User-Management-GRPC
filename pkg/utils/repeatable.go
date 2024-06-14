package utils

import "time"

func DoWithTries(fnc func() error, tries int, delay time.Duration) (err error) {
	for tries > 0 {
		if err = fnc(); err == nil {
			time.Sleep(delay)
			tries--

			continue
		}

		return nil
	}

	return
}
