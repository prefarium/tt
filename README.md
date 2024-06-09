I made a very simple tool for tracking elapsed work time for myself

It requires some configuration before using as follows:
```bash
go build -ldflags="-X 'main.csvPath=/Users/user/db.csv' -X 'main.offset=3h' -X 'main.location=Europe/Moscow'"
```
* main.csvPath - path to a file for storing tracked time, you must create it yourself
* main.offset - affects what time your day starts and ends for your current day time tracking;
  meaning that offset of 3h makes your days start and end at 3 am
* main.location - sets time zone for week time calculation

After compilation simply execute the binary to start tracking today's time.
Press CTRL+C to stop tracking and write worked time to CSV.
If you run the binary again it will consider today's previous time

To see the time you worked during the whole week simply use `tt week`
