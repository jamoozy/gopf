<?php
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
