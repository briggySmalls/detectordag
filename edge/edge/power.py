from edge.ina219 import INA219


class Power:
    def __init__(self, ina219: INA219):
        self.ina219 = ina219

    def is_powered(self) -> bool:
        current = self.ina219.getCurrent_mA()
        return current > 0
