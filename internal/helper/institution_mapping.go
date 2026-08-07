package helper

var UserInstitutionMap = map[uint][]uint{
    94: {2},
    // 2: {1, 2},
    // 3: {2},
}

func HasInstitutionAccess(userID uint, institutionID uint) bool {
    institutions, ok := UserInstitutionMap[userID]
    if !ok {
        return false
    }

    for _, id := range institutions {
        if id == institutionID {
            return true
        }
    }

    return false
}