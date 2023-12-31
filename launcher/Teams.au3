#NoTrayIcon
$file = StringReplace(@ScriptName, '.exe', '.ini')
If FileExists($file) Then
	$iniData = IniReadSection($file, 'data')
	If Not @error Then
		$data = $iniData[1][1]
	Else
		$data = 'data'
	EndIf
	DirCreate(@ScriptDir & "\" & $data & "\Users\Current")
	DirCreate(@ScriptDir & "\" & $data & "\Users\Current\Temp")
	DirCreate(@ScriptDir & "\" & $data & "\Users\Current\AppData\Local")
	DirCreate(@ScriptDir & "\" & $data & "\Users\Current\AppData\Roaming")
	DirCreate(@ScriptDir & "\" & $data & "\Users\Public")
	DirCreate(@ScriptDir & "\" & $data & "\ProgramData")
	EnvSet("Temp", @ScriptDir & "\" & $data & "\Users\Current\Temp")
	EnvSet("Tmp", @ScriptDir & "\" & $data & "\Users\Current\Temp")
	EnvSet("AppData", @ScriptDir & "\" & $data & "\Users\Current\AppData\Roaming")
	EnvSet("LocalAppData", @ScriptDir & "\" & $data & "\Users\Current\AppData\Local")
	EnvSet("UserProfile", @ScriptDir & "\" & $data & "\Users\Current")
	EnvSet("Public", @ScriptDir & "\" & $data & "\Users\Public")
	EnvSet("ProgramData", @ScriptDir & "\" & $data & "\ProgramData")
	EnvSet("AllUsersProfile", @ScriptDir & "\" & $data & "\ProgramData")
	EnvSet("HomeDrive", @ScriptDir & "\" & $data)
	EnvSet("HomePath", "\Users\Current")
	EnvSet("SystemDrive", @ScriptDir & "\" & $data)
	$iniEnv = IniReadSection($file, 'env')
	If Not @error Then
		For $i = 0 To $iniEnv[0][0]
			EnvSet($iniEnv[$i][0], EnvGet($iniEnv[$i][0]) & ";" & $iniEnv[$i][1])
		Next
	EndIf
	$iniStart = IniReadSection($file, 'start')
	If Not @error Then
		For $i = 0 To $iniStart[0][0]
			If StringIsInt($iniStart[$i][0]) Then
				Run(@ScriptDir & '\' & $iniStart[$i][1], @ScriptDir)
			Else
				Run(@ScriptDir & '\' & $iniStart[$i][1], @ScriptDir & '\' & $iniStart[$i][0])
			EndIf
		Next
	Else
		$nFile = StringReplace(@ScriptName, 'Launcher', '')
		Run(@ScriptDir & '\' & $nFile, @ScriptDir)
	EndIf
Else
	FileWriteLine($file, "[start]")
	FileWriteLine($file, ";设置要启动的程序，可以带启动参数")
	FileWriteLine($file, ";等号前面的内容为启动引用的位置，数字则表示使用当前目录")
	FileWriteLine($file, ";0=")
	FileWriteLine($file, "")
	FileWriteLine($file, "[data]")
	FileWriteLine($file, ";设置数据文件保存的文件夹名")
	FileWriteLine($file, ";相对于当前目录，仅读取第一个参数")
	FileWriteLine($file, ";0=")
	FileWriteLine($file, "")
	FileWriteLine($file, "[env]")
	FileWriteLine($file, ";设置环境变量")
	FileWriteLine($file, ";按 Windows 环境变量格式设置，已存在的变量将添加至尾部")
	FileWriteLine($file, ";PATH=.")
EndIf