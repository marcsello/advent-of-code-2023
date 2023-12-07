{$APPTYPE CONSOLE}
program day6_part2;

const
  Time:QWord = 56977793;
  Distance:QWord = 499221010971440;


var
  holdTime: QWord;
  possibilities: QWord;
  tempVar: QWord; // Ah yes, I have fond memories of this convention of my Pascal days


begin

    possibilities := 0;

    for holdTime := 0 to Time do
    begin
      tempVar := (Time - holdTime) * holdTime;
      if tempVar > Distance then
      begin
        possibilities += 1;
      end;
    end;

   WriteLn(possibilities);
end.

