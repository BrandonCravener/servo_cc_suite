#include <Servo.h>

// Config
int servoPins[] = {2, 3, 4, 5, 6, 7};

// DO NOT TOUCH ANYTHING BELOW HERE
Servo servos[6];

void setup() {
  Serial.begin(115200);
  for (int i = 0; i < 6; i++) {
    int servoPin = servoPins[i];
    servos[i].attach(servoPin);
    delay(50);
    servos[i].write(90);
    delay(416);
  }
}
void loop() {
}
