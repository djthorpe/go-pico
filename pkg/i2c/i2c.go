package i2c



i2c := machine.I2C0
err := i2c.Configure(machine.I2CConfig{
	SCL: machine.P0_30,
	SDA: machine.P0_31,
})
if err != nil {
	println("could not configure I2C:", err)
	return
}

w := []byte{0x75}
r := make([]byte, 1)
err = i2c.Tx(0x68, w, r)
if err != nil {
	println("could not interact with I2C device:", err)
	return
}
println("WHO_AM_I:", r[0]) // prints "WHO_AM_I: 104"
