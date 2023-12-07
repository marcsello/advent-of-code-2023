{$APPTYPE CONSOLE}
program day6_part1;

const
  Times:     Array[0..3] of integer = (56, 97, 77, 93);
  Distances: Array[0..3] of integer = (499, 2210, 1097, 1440);


var
  i: byte;
  holdTime: integer;
  possibilities:Array[0..3] of integer;
  tempVar: integer; // Ah yes, I have fond memories of this convention of my Pascal days


begin
   For i := 0 to 3 do
   begin
      possibilities[i] := 0;

      for holdTime := 0 to Times[i] do
      begin
        tempVar := (Times[i] - holdTime) * holdTime;
        if tempVar > Distances[i] then
        begin
          possibilities[i] += 1;
        end;
      end;

   end;

   tempVar := possibilities[0];
   for i := 1 to 3 do
   begin
     tempVar *= possibilities[i]
   end;

   WriteLn(tempVar);
end.

