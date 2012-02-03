<?php
// Copyright 2012 Andrew "Jamoozy" Correa
//
// This file is part of GOPF.
//
// GOPF is free software: you can redistribute it and/or modify it under
// the terms of the GNU General Public as published by the Free Software
// Foundation, either version 3 of the License, or (at your option) any
// later version.
//
// GOPF is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License
// for more details.
//
// You should have received a copy of the GNU General Public License
// along with GOPF. If not, see http://www.gnu.org/licenses/.

$log_file = 'logged_ips.log';

// Determines if the IP address is bad.  Returns an error log to present 
// to the user if the IP is bad, and false otherwise.
function ip_is_ok($ip) {
  global $log_file;

  // Write "ip, date" to log file.
  $line = "$ip  ".date("Y-m-d H:i:s")."\n";

  if ($file = fopen($log_file, "a")) {
    if (!!fwrite($file, $line)) {
      $rtn = fclose($file);
    } else {
      $rtn = "Could not close $log_file<br />for $line";
    }

    if (!$rtn) {
      return "Logged: $line";
    } else {
      return false;
    }
  } else {
    return "Could not open $log_file<br />for $line";
  }
}
?>
