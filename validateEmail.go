package main

import "github.com/badoux/checkmail"

func IsValidEmail(email string) error {
    if err := checkmail.ValidateFormat(email); err != nil {
        return err
    }

    if err := checkmail.ValidateHost(email); err != nil {
        return err
    }

    return nil
}