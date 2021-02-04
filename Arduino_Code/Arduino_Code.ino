#include <Servo.h>

// Config
#define COMM_SIZE   1
#define INIT_DELAY  750
#define DEG_OFFSET  30
int servoPins[] = {2, 3, 4, 5, 6, 7}; // e, a, d, g, b, e2

// DO NOT TOUCH ANYTHING BELOW HERE
Servo servo_e;
Servo servo_a;
Servo servo_d;
Servo servo_g;
Servo servo_b;
Servo servo_e_2;

Servo servos[] = {servo_e, servo_a, servo_d, servo_g, servo_b, servo_e_2};

void setup() {
  // put your setup code here, to run once:
  Serial.begin(115200);
  for (int i = 0; i < 6; i++) {
    int servoPin = servoPins[i];
    servos[i].attach(servoPin);
    delay(50);
  }
  for (int i = 0; i < 6; i++) {
    servos[i].write(90 + DEG_OFFSET);
  }
  delay(500);
  Serial.print("DONE");
}

void loop() {
  if (Serial.available() > 0) {
    char input = Serial.read();
    //    Serial.println(input);
    int servoID = atoi(&input);
    if (servoID < 6 && servoID >= 0) {
      if (servos[servoID].read() < 90) {
        servos[servoID].write(90 + DEG_OFFSET);
      } else {
        servos[servoID].write(90 - DEG_OFFSET);
      }
    }
    while (Serial.available()) {
      Serial.read();
    }
  }
}
