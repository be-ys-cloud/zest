package structures

type Statistics struct {
	SavedBytes      int `sql:"savedBytes"`
	NbFiles         int `sql:"nbFiles"`
	TotalSizeOnDisk int `sql:"totalSizeOnDisk"`
}
