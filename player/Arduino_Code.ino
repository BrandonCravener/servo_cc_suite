#include <Servo.h>

// Config
int servoPins[] = {2, 3, 4, 5, 6, 7};
int startingDEG[] = {80, 80, 80, 80, 80, 80};
int endingDEG[] = {100, 100, 100, 100, 100, 100};

// DO NOT TOUCH ANYTHING BELOW HERE
Servo servos[6];

void setup() {
  Serial.begin(115200);
  for (int i = 0; i < 6; i++) {
    int servoPin = servoPins[i];
    servos[i].attach(servoPin);
    delay(50);
    servos[i].write(startingDEG[i]);
    delay(416);
  }
  Serial.println("DONE");
}
void loop() {
  if (Serial.available() > 0) {
    char input = Serial.read();
    int servoID = atoi(&input);
    if (servoID < 6 && servoID >= 0) {
      if (servos[servoID].read() <= startingDEG[servoID] + 3) {
        servos[servoID].write(endingDEG[servoID]);
      } else {
        servos[servoID].write(startingDEG[servoID]);
      }
    }
    while (Serial.available()) {
      Serial.read();
    }
  }
}
